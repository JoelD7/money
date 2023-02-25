package mocks

import (
	"context"
	"errors"
	"github.com/JoelD7/money/api/storage"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ForceNotFound   = false
	ForceUserExists = false

	ErrForceNotFound = errors.New("force not found")
)

type MockDynamo struct{}

func InitDynamoMock() {
	storage.DefaultClient = &MockDynamo{}
}

func (d *MockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if ForceNotFound {
		return &dynamodb.GetItemOutput{}, ErrForceNotFound
	}

	email, err := attributevalue.Marshal("test@gmail.com")
	if err != nil {
		return nil, err
	}

	password, err := attributevalue.Marshal("$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe")
	if err != nil {
		return nil, err
	}

	fullName, err := attributevalue.Marshal("Joel")
	if err != nil {
		return nil, err
	}

	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"email":     email,
			"password":  password,
			"full_name": fullName,
		},
	}, nil
}

func (d *MockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if ForceNotFound {
		return &dynamodb.QueryOutput{}, ErrForceNotFound
	}

	email, err := attributevalue.Marshal("test@gmail.com")
	if err != nil {
		return nil, err
	}

	password, err := attributevalue.Marshal("$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe")
	if err != nil {
		return nil, err
	}

	fullName, err := attributevalue.Marshal("Joel")
	if err != nil {
		return nil, err
	}

	return &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				"email":     email,
				"password":  password,
				"full_name": fullName,
			},
		},
	}, nil
}

func (d *MockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if ForceUserExists {
		return &dynamodb.PutItemOutput{}, storage.ErrExistingUser
	}

	return &dynamodb.PutItemOutput{}, nil
}
