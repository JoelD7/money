package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/entities"
	"github.com/JoelD7/money/api/shared/env"
	"github.com/JoelD7/money/api/shared/utils"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	tableName = env.GetString("USERS_TABLE_NAME", "person")

	errNotFound     = errors.New("person not found")
	ErrExistingUser = errors.New("this account already exists")
)

const (
	categoryPrefix = "CTG"
	personPrefix   = "PS"
)

var (
	emailIndex = "email-index"
)

func CreatePerson(ctx context.Context, fullName, email, password string) error {
	person := &entities.Person{
		PersonID:   utils.GenerateDynamoID(personPrefix),
		FullName:   fullName,
		Email:      email,
		Password:   password,
		Categories: getDefaultCategories(),
	}

	item, err := attributevalue.MarshalMap(person)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	}

	_, err = DefaultClient.PutItem(ctx, input)
	if err != nil {
		fmt.Println("storage: ", err)
		return ErrExistingUser
	}

	if err != nil {
		return err
	}

	return nil
}

func GetPerson(ctx context.Context, personId string) (*entities.Person, error) {
	personKey, err := attributevalue.Marshal(personId)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"person_id": personKey,
		},
	}

	result, err := DefaultClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errNotFound
	}

	person := new(entities.Person)
	err = attributevalue.UnmarshalMap(result.Item, person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func GetPersonByEmail(ctx context.Context, email string) (*entities.Person, error) {
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
		TableName:                 aws.String(tableName),
	}

	result, err := DefaultClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, errNotFound
	}

	person := new(entities.Person)
	err = attributevalue.UnmarshalMap(result.Items[0], person)
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
