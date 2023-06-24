package expenses

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

type DynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

const (
	splitter = ":"
)

var (
	dynamoClient  *dynamodb.Client
	DefaultClient DynamoAPI

	awsRegion = env.GetString("REGION", "us-east-1")

	TableName                = env.GetString("EXPENSES_TABLE_NAME", "expenses")
	periodUserExpenseIDIndex = "period_user-expense_id-index"

	ErrNotFound    = errors.New("expenses not found")
	ErrEmptyUserID = errors.New("empty userID")
	ErrEmptyPeriod = errors.New("empty period")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func GetExpensesByPeriod(ctx context.Context, userID, periodID string) ([]*models.Expense, error) {
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

	result, err := DefaultClient.Query(ctx, input)
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
