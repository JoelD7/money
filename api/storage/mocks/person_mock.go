package mocks

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/storage"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ForceNotFound   = false
	ForceUserExists = false

	ErrForceNotFound      = errors.New("force not found")
	ErrMockNotInitialized = errors.New("mock is not initialized")
)

type MockDynamo struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput
}

func InitDynamoMock() *MockDynamo {
	getItemOutput, queryOutput, err := defaultOutput()
	if err != nil {
		panic(fmt.Errorf("initDynamoMock: %w", err))
	}

	mock := &MockDynamo{
		GetItemOutput: getItemOutput,
		QueryOutput:   queryOutput,
	}

	storage.DefaultClient = mock

	return mock
}

func (d *MockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if ForceNotFound {
		return &dynamodb.GetItemOutput{}, ErrForceNotFound
	}

	if d.GetItemOutput == nil {
		return &dynamodb.GetItemOutput{}, ErrMockNotInitialized
	}

	return d.GetItemOutput, nil
}

func (d *MockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if ForceNotFound {
		return &dynamodb.QueryOutput{}, ErrForceNotFound
	}

	if d.QueryOutput == nil {
		return &dynamodb.QueryOutput{}, ErrMockNotInitialized
	}

	return d.QueryOutput, nil
}

func (d *MockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if ForceUserExists {
		return &dynamodb.PutItemOutput{}, storage.ErrExistingUser
	}

	return &dynamodb.PutItemOutput{}, nil
}

func defaultOutput() (*dynamodb.GetItemOutput, *dynamodb.QueryOutput, error) {
	email, err := attributevalue.Marshal("test@gmail.com")
	if err != nil {
		return nil, nil, err
	}

	password, err := attributevalue.Marshal("$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe")
	if err != nil {
		return nil, nil, err
	}

	fullName, err := attributevalue.Marshal("Joel")
	if err != nil {
		return nil, nil, err
	}

	return &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"email":     email,
				"password":  password,
				"full_name": fullName,
			},
		},
		&dynamodb.QueryOutput{
			Items: []map[string]types.AttributeValue{
				{
					"email":     email,
					"password":  password,
					"full_name": fullName,
				},
			},
		},
		nil
}
