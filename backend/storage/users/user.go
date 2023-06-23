package users

import (
	"context"
	"errors"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	TableName = env.GetString("USERS_TABLE_NAME", "users")

	ErrNotFound     = errors.New("user not found")
	ErrExistingUser = errors.New("this account already exists")

	emailIndex = "email-index"
)

const (
	categoryPrefix = "CTG"
	userPrefix     = "US"
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func CreateUser(ctx context.Context, fullName, email, password string) error {
	ok, err := userExists(ctx, email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if ok {
		return ErrExistingUser
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

	_, err = DefaultClient.PutItem(ctx, input)
	if err != nil {
		return ErrExistingUser
	}

	if err != nil {
		return err
	}

	return nil
}

func userExists(ctx context.Context, email string) (bool, error) {
	user, err := GetUserByEmail(ctx, email)
	if user != nil {
		return true, nil
	}

	return false, err
}

func GetUser(ctx context.Context, userID string) (*models.User, error) {
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

	result, err := DefaultClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	user := new(models.User)
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
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

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	user := new(models.User)
	err = attributevalue.UnmarshalMap(result.Items[0], user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(ctx context.Context, user *models.User) error {
	updatedItem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedItem,
		TableName: aws.String(TableName),
	}

	_, err = DefaultClient.PutItem(ctx, input)
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
