package users

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"time"
)

var (
	TableName = env.GetString("USERS_TABLE_NAME", "users")

	emailIndex = "email-index"
)

const (
	categoryPrefix = "CTG"
	userPrefix     = "US"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) createUser(ctx context.Context, fullName, email, password string) error {
	ok, err := d.userExists(ctx, email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return err
	}

	if ok {
		return models.ErrExistingUser
	}

	user := &models.User{
		UserID:      utils.GenerateDynamoID(userPrefix),
		FullName:    fullName,
		Email:       email,
		Password:    password,
		Categories:  getDefaultCategories(),
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return models.ErrExistingUser
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoRepository) userExists(ctx context.Context, email string) (bool, error) {
	user, err := d.getUserByEmail(ctx, email)
	if user != nil {
		return true, nil
	}

	return false, err
}

func (d *DynamoRepository) getUser(ctx context.Context, userID string) (*models.User, error) {
	userKey, err := attributevalue.Marshal(userID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"user_id": userKey,
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, models.ErrUserNotFound
	}

	user := new(models.User)
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DynamoRepository) getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	nameEx := expression.Name("email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		IndexName:                 aws.String(emailIndex),
		TableName:                 aws.String(TableName),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrUserNotFound
	}

	user := new(models.User)
	err = attributevalue.UnmarshalMap(result.Items[0], user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DynamoRepository) updateUser(ctx context.Context, user *models.User) error {
	updatedItem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedItem,
		TableName: aws.String(TableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	return err
}

func getDefaultCategories() []*models.Category {
	return []*models.Category{
		{
			CategoryID:   utils.GenerateDynamoID(categoryPrefix),
			CategoryName: "Entertainment",
			Color:        "#ff8733",
		},
		{
			CategoryID:   utils.GenerateDynamoID(categoryPrefix),
			CategoryName: "Health",
			Color:        "#00b85e",
		},
		{
			CategoryID:   utils.GenerateDynamoID(categoryPrefix),
			CategoryName: "Utilities",
			Color:        "#009eb8",
		},
	}
}
