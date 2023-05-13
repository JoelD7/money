package invalid_token

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
	"time"
)

type TokenType string

type DynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

var (
	InvalidTokenTableName = env.GetString("INVALID_TOKEN_TABLE_NAME", "invalid_token")

	errNotFound = errors.New("no tokens found for this user")
)

var (
	dynamoClient  *dynamodb.Client
	DefaultClient DynamoAPI

	awsRegion = env.GetString("REGION", "us-east-1")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func AddInvalidToken(ctx context.Context, email, token string, tokenType TokenType, expires int64) error {
	invalidToken := models.InvalidToken{
		Email:       email,
		Token:       token,
		Expire:      expires,
		Type:        string(tokenType),
		CreatedDate: time.Now(),
	}

	item, err := attributevalue.MarshalMap(invalidToken)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(InvalidTokenTableName),
	}

	_, err = DefaultClient.PutItem(ctx, input)

	return err
}

func GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error) {
	nameEx := expression.Name("email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		TableName:                 aws.String(InvalidTokenTableName),
	}

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, errNotFound
	}

	invalidTokens := make([]*models.InvalidToken, 0)

	err = attributevalue.UnmarshalListOfMaps(result.Items, invalidTokens)
	if err != nil {
		return nil, err
	}

	return invalidTokens, nil
}
