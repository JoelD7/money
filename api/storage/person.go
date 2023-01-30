package storage

import (
	"errors"
	"github.com/JoelD7/money/api/entities"
	"github.com/JoelD7/money/api/shared/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	db        *dynamodb.DynamoDB
	tableName = env.GetString("USERS_TABLE_NAME", "person")
	awsRegion = env.GetString("REGION", "us-east-1")

	errNotFound = errors.New("person not found")
)

func init() {
	dynamodbSession, err := session.NewSession(aws.NewConfig().WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	db = dynamodb.New(dynamodbSession)
}

func GetPerson(personId string) (*entities.Person, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"person_id": {
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
