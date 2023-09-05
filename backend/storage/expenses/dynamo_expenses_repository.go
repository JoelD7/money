package expenses

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	splitter = ":"
)

var (
	TableName                = env.GetString("EXPENSES_TABLE_NAME", "expenses")
	periodUserExpenseIDIndex = "period_user-expense_id-index"

	ErrNotFound    = errors.New("expenses not found")
	ErrEmptyUserID = errors.New("empty userID")
	ErrEmptyPeriod = errors.New("empty period")
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) GetExpensesByPeriod(ctx context.Context, userID, periodID string) ([]*models.Expense, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}

	if periodID == "" {
		return nil, ErrEmptyPeriod
	}

	periodUser := periodID + splitter + userID

	nameEx := expression.Name("period_user").Equal(expression.Value(periodUser))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		IndexName:                 aws.String(periodUserExpenseIDIndex),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	expenses := make([]*models.Expense, 0)
	err = attributevalue.UnmarshalListOfMaps(result.Items, &expenses)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
