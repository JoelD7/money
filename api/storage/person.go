package storage

import (
	"errors"
	"github.com/JoelD7/money/api/entities"
	"github.com/JoelD7/money/api/shared/env"
	"github.com/JoelD7/money/api/shared/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strings"
)

var (
	db        *dynamodb.DynamoDB
	tableName = env.GetString("USERS_TABLE_NAME", "person")
	awsRegion = env.GetString("REGION", "us-east-1")

	errNotFound     = errors.New("person not found")
	ErrExistingUser = errors.New("this account already exists")
)

const (
	categoryPrefix = "CTG"
)

func init() {
	dynamodbSession, err := session.NewSession(aws.NewConfig().WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	db = dynamodb.New(dynamodbSession)
}

func CreatePerson(fullName, email, password string) error {
	person := &entities.Person{
		FullName:   fullName,
		Email:      email,
		Password:   password,
		Categories: getDefaultCategories(),
	}

	item, err := dynamodbattribute.MarshalMap(person)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	}

	_, err = db.PutItem(input)
	if err != nil && strings.Contains(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
		return ErrExistingUser
	}

	if err != nil {
		return err
	}

	return nil
}

func GetPerson(personId string) (*entities.Person, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(personId),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errNotFound
	}

	person := new(entities.Person)
	err = dynamodbattribute.UnmarshalMap(result.Item, person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func getDefaultCategories() []*entities.Category {
	return []*entities.Category{
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
