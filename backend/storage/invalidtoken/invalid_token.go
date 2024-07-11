package invalidtoken

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

type Type string

type DynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

const (
	TypeAccess  Type = "access"
	TypeRefresh Type = "refresh"
)

var (
	tableName = env.GetString("INVALID_TOKEN_TABLE_NAME", "")

	ErrNotFound = errors.New("no tokens found for this user")
)

var (
	dynamoClient  *dynamodb.Client
	DefaultClient DynamoAPI

	awsRegion = env.GetString("AWS_REGION", "")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func Add(ctx context.Context, email, token string, tokenType Type, expires int64) error {
	invalidToken := models.InvalidToken{
		Token:       token,
		Expire:      expires,
		CreatedDate: time.Now(),
	}

	item, err := attributevalue.MarshalMap(invalidToken)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = DefaultClient.PutItem(ctx, input)

	return err
}

func GetAllForPerson(ctx context.Context, email string) ([]models.InvalidToken, error) {
	nameCondition := expression.Name("email").Equal(expression.Value(email))
	filterCondition := expression.Name("expire").GreaterThanEqual(expression.Value(time.Now().Unix()))

	expr, err := expression.NewBuilder().
		WithCondition(nameCondition).
		WithFilter(filterCondition).
		Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.Condition(),
		TableName:                 aws.String(tableName),
	}

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	invalidTokens := make([]models.InvalidToken, 0)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &invalidTokens)
	if err != nil {
		return nil, err
	}

	return invalidTokens, nil
}
