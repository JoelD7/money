package person

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"time"
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

	UsersTableName = env.GetString("USERS_TABLE_NAME", "person")

	ErrNotFound     = errors.New("person not found")
	ErrExistingUser = errors.New("this account already exists")

	emailIndex = "email-index"
)

const (
	categoryPrefix = "CTG"
	personPrefix   = "PS"
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	DefaultClient = dynamoClient
}

func CreatePerson(ctx context.Context, fullName, email, password string) error {
	ok, err := personExists(ctx, email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if ok {
		return ErrExistingUser
	}

	person := &models.Person{
		PersonID:    utils.GenerateDynamoID(personPrefix),
		FullName:    fullName,
		Email:       email,
		Password:    password,
		Categories:  getDefaultCategories(),
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}

	item, err := attributevalue.MarshalMap(person)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(UsersTableName),
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

func personExists(ctx context.Context, email string) (bool, error) {
	person, err := GetPersonByEmail(ctx, email)
	if person != nil {
		return true, nil
	}

	return false, err
}

func GetPerson(ctx context.Context, personId string) (*models.Person, error) {
	personKey, err := attributevalue.Marshal(personId)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]types.AttributeValue{
			"person_id": personKey,
		},
	}

	result, err := DefaultClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	person := new(models.Person)
	err = attributevalue.UnmarshalMap(result.Item, person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func GetPersonByEmail(ctx context.Context, email string) (*models.Person, error) {
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
		TableName:                 aws.String(UsersTableName),
	}

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	person := new(models.Person)
	err = attributevalue.UnmarshalMap(result.Items[0], person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func UpdatePerson(ctx context.Context, person *models.Person) error {
	updatedItem, err := attributevalue.MarshalMap(person)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      updatedItem,
		TableName: aws.String(UsersTableName),
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
