package users

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
	"time"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
	tableName    string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, tableName string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}
	tableNameEnv := env.GetString("USERS_TABLE_NAME", "")

	if tableNameEnv == "" && tableName == "" {
		return nil, fmt.Errorf("initialize income recurring dynamo repository failed: table name is required")
	}

	d.tableName = tableName
	if d.tableName == "" {
		d.tableName = tableNameEnv
	}

	return d, nil
}

func (d *DynamoRepository) CreateUser(ctx context.Context, u *models.User) error {
	user := toUserEntity(u)

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"),
		TableName:           aws.String(d.tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrExistingUser)
	}

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	userKey, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username": userKey,
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, models.ErrUserNotFound
	}

	user := new(userEntity)
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return toUserModel(user), nil
}

func (d *DynamoRepository) UpdateUser(ctx context.Context, u *models.User) error {
	user := toUserEntity(u)

	user.UpdatedDate = time.Now()

	updatedItem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedItem,
		TableName: aws.String(d.tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	return err
}

func (d *DynamoRepository) DeleteUser(ctx context.Context, username string) error {
	userKey, err := attributevalue.Marshal(username)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username": userKey,
		},
	}

	_, err = d.dynamoClient.DeleteItem(ctx, input)
	return err
}
