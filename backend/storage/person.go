package storage

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/entities"
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
	ok, err := personExists(ctx, email)
	if err != nil && !errors.Is(err, errNotFound) {
		return err
	}

	if ok {
		return ErrExistingUser
	}

	person := &entities.Person{
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
		TableName: aws.String(tableName),
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

func UpdatePerson(ctx context.Context, person *entities.Person) error {
	updatedItem, err := attributevalue.MarshalMap(person)
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"person_id": &types.AttributeValueMemberS{Value: person.PersonID},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#checks_settings": "checks_settings",
		},
		ExpressionAttributeValues: updatedItem,
		UpdateExpression:          aws.String("SET #person = :person"),
	}

	_, err = DefaultClient.UpdateItem(ctx, input)
	return err
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
