package income

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

var (
	dynamoClient  *dynamodb.Client
	DefaultClient DynamoAPI

	awsRegion = env.GetString("REGION", "us-east-1")

	TableName         = env.GetString("INCOME_TABLE_NAME", "income")
	userIdPeriodIndex = "user_id-period_id-index"

	ErrNotFound = errors.New("income not found")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func GetIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error) {
	userIDEx := expression.Name("user_id").Equal(expression.Value(userID))
	periodEx := expression.Name("period_id").Equal(expression.Value(periodID))
	nameEx := userIDEx.And(periodEx)

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		IndexName:                 aws.String(userIdPeriodIndex),
	}

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	income := make([]*models.Income, 0)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &income)
	if err != nil {
		return nil, err
	}

	return income, nil
}
