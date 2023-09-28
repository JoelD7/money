package users

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"time"
)

var (
	TableName = env.GetString("USERS_TABLE_NAME", "users")
)

const (
	categoryPrefix = "CTG"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) CreateUser(ctx context.Context, fullName, username, password string) error {
	user := &models.User{
		FullName:    fullName,
		Username:    username,
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
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"),
		TableName:           aws.String(TableName),
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
		TableName: aws.String(TableName),
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

	user := new(models.User)
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DynamoRepository) UpdateUser(ctx context.Context, user *models.User) error {
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
			ID:    utils.GenerateDynamoID(categoryPrefix),
			Name:  aws.String("Entertainment"),
			Color: aws.String("#ff8733"),
		},
		{
			ID:    utils.GenerateDynamoID(categoryPrefix),
			Name:  aws.String("Health"),
			Color: aws.String("#00b85e"),
		},
		{
			ID:    utils.GenerateDynamoID(categoryPrefix),
			Name:  aws.String("Utilities"),
			Color: aws.String("#009eb8"),
		},
	}
}
