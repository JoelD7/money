package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	db        = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	tableName = "person"

	errNotFound = errors.New("person not found")
)

func getItem(personId string) (*Person, error) {
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

	person := new(Person)
	err = dynamodbattribute.UnmarshalMap(result.Item, person)
	if err != nil {
		return nil, err
	}

	return person, nil
}
