package invalid_token

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	InvalidTokenTableName = env.GetString("INVALID_TOKEN_TABLE_NAME", "invalid_token")

	errNotFound = errors.New("no tokens found for this user")
)

func AddInvalidToken(ctx context.Context, email, token string, expires int64) error {
	invalidToken := models.InvalidToken{
		Email:  email,
		Token:  token,
		Expire: expires,
	}

	item, err := attributevalue.MarshalMap(invalidToken)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(InvalidTokenTableName),
	}

	_, err = storage.DefaultClient.PutItem(ctx, input)

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

	result, err := storage.DefaultClient.Query(ctx, input)
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
