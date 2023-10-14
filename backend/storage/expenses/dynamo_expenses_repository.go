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
	"strings"
	"time"
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

func (d *DynamoRepository) CreateExpense(ctx context.Context, expense *models.Expense) error {
	entity := toExpenseEntity(expense)

	entity.CreatedDate = time.Now()
	entity.PeriodUser = buildPeriodUser(entity.Username, entity.Period)

	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return fmt.Errorf("marshal expense failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("put expense failed: %v", err)
	}

	return nil
}

func (d *DynamoRepository) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	entity := toExpenseEntity(expense)

	entity.UpdateDate = time.Now()

	username, err := attributevalue.Marshal(entity.Username)
	if err != nil {
		return fmt.Errorf("marshaling username key: %v", err)
	}

	expenseID, err := attributevalue.Marshal(entity.ExpenseID)
	if err != nil {
		return fmt.Errorf("marshaling expense id key: %v", err)
	}

	attributeValues, err := getAttributeValues(entity)
	if err != nil {
		return fmt.Errorf("get attribute values failed: %v", err)
	}

	updateExpression := getUpdateExpression(attributeValues)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"username":   username,
			"expense_id": expenseID,
		},
		TableName:                 aws.String(tableName),
		ConditionExpression:       aws.String("attribute_exists(expense_id)"),
		ExpressionAttributeValues: attributeValues,
		UpdateExpression:          updateExpression,
	}

	_, err = d.dynamoClient.UpdateItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrExpensesNotFound)
	}

	if err != nil {
		return fmt.Errorf("put expense failed: %v", err)
	}

	return nil
}

func getAttributeValues(expense *expenseEntity) (map[string]types.AttributeValue, error) {
	attrValues := make(map[string]types.AttributeValue)

	categoryID, err := attributevalue.Marshal(expense.CategoryID)
	if err != nil {
		return nil, err
	}

	amount, err := attributevalue.Marshal(expense.Amount)
	if err != nil {
		return nil, err
	}

	name, err := attributevalue.Marshal(expense.Name)
	if err != nil {
		return nil, err
	}

	notes, err := attributevalue.Marshal(expense.Notes)
	if err != nil {
		return nil, err
	}

	updatedDate, err := attributevalue.Marshal(time.Now())
	if err != nil {
		return nil, err
	}

	if expense.CategoryID != nil {
		attrValues[":category_id"] = categoryID
	}

	if expense.Amount != 0 {
		attrValues[":amount"] = amount
	}

	if expense.Name != "" {
		attrValues[":name"] = name
	}

	if expense.Notes != "" {
		attrValues[":notes"] = notes
	}

	attrValues[":update_date"] = updatedDate

	return attrValues, nil
}

func getUpdateExpression(attributeValues map[string]types.AttributeValue) *string {
	attributes := make([]string, 0)

	if _, ok := attributeValues[":category_id"]; ok {
		attributes = append(attributes, "category_id = :category_id")
	}

	if _, ok := attributeValues[":amount"]; ok {
		attributes = append(attributes, "amount = :amount")
	}

	if _, ok := attributeValues[":name"]; ok {
		attributes = append(attributes, "name = :name")
	}

	if _, ok := attributeValues[":notes"]; ok {
		attributes = append(attributes, "notes = :notes")
	}

	if _, ok := attributeValues[":update_date"]; ok {
		attributes = append(attributes, "update_date = :update_date")
	}

	return aws.String("SET " + strings.Join(attributes, ", "))
}

func (d *DynamoRepository) GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"username":   &types.AttributeValueMemberS{Value: username},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get expense failed: %v", err)
	}

	if result.Item == nil {
		return nil, models.ErrExpenseNotFound
	}

	entity := new(expenseEntity)

	err = attributevalue.UnmarshalMap(result.Item, entity)
	if err != nil {
		return nil, fmt.Errorf("unmarshal expense item failed: %v", err)
	}

	return toExpenseModel(entity), nil
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

func (d *DynamoRepository) DeleteExpense(ctx context.Context, expenseID, username string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"username":   &types.AttributeValueMemberS{Value: username},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	}

	_, err := d.dynamoClient.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("delete expense failed: %v", err)
	}

	return nil
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
