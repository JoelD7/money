package expenses

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	defaultPageSize = 10
)

var (
	tableName                = env.GetString("EXPENSES_TABLE_NAME", "expenses")
	periodUserExpenseIDIndex = "period_user-expense_id-index"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, "", startKey, nil, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input)
}

func (d *DynamoRepository) GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, periodID, startKey, categories, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input)
}

func (d *DynamoRepository) GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, periodID, startKey, nil, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input)
}

func (d *DynamoRepository) GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, "", startKey, categories, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input)
}

func buildQueryInput(username, periodID, startKey string, categories []string, pageSize int) (*dynamodb.QueryInput, error) {
	conditionEx := expression.Name("username").Equal(expression.Value(username))

	var decodedStartKey map[string]types.AttributeValue
	var err error

	if startKey != "" {
		decodedStartKey, err = decodeStartKey(startKey)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	input := &dynamodb.QueryInput{
		TableName:         aws.String(tableName),
		ExclusiveStartKey: decodedStartKey,
		Limit:             getPageSize(pageSize),
	}

	// Query the period_user-expense_id-index
	if periodID != "" {
		input.IndexName = aws.String(periodUserExpenseIDIndex)

		periodUser := buildPeriodUser(username, periodID)
		conditionEx = expression.Name("period_user").Equal(expression.Value(periodUser))
	}

	conditionBuilder := expression.NewBuilder().WithCondition(conditionEx)

	if categories == nil || len(categories) > 0 {
		filterCondition := buildCategoriesConditionFilter(categories)
		conditionBuilder = conditionBuilder.WithFilter(filterCondition)
	}

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, err
	}

	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()
	input.KeyConditionExpression = expr.Condition()

	return input, nil
}

func buildCategoriesConditionFilter(categories []string) expression.ConditionBuilder {
	conditions := make([]expression.ConditionBuilder, 0, len(categories))

	for _, categoryID := range categories {
		conditions = append(conditions, expression.Name("category_id").Equal(expression.Value(categoryID)))
	}

	if len(categories) == 1 {
		return expression.Name("category_id").Equal(expression.Value(categories[0]))
	}

	if len(categories) == 2 {
		return expression.Or(conditions[0], conditions[1])
	}

	return expression.Or(conditions[0], conditions[1], conditions[2:]...)
}

func (d *DynamoRepository) performQuery(ctx context.Context, input *dynamodb.QueryInput) ([]*models.Expense, string, error) {
	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrExpensesNotFound
	}

	expensesEntities := new([]*expenseEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &expensesEntities)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal expenses items failed: %v", err)
	}

	nextKey, err := encodeLastKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toExpenseModels(*expensesEntities), nextKey, nil
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}

func buildPeriodUser(username, period string) string {
	return fmt.Sprintf("%s:%s", period, username)
}
