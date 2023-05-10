package person

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/storage"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FailureCondition string

const (
	NotFound   FailureCondition = "not found"
	UserExists FailureCondition = "user exists"
	None       FailureCondition = "none"
)

var (
	ErrForceNotFound      = errors.New("force not found")
	ErrMockNotInitialized = errors.New("mock is not initialized")
)

type MockDynamo struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput

	emulatingErrors map[FailureCondition]error
	mockedErr       error
}

func InitDynamoMock() *MockDynamo {
	getItemOutput, queryOutput, err := defaultOutput()
	if err != nil {
		panic(fmt.Errorf("initDynamoMock: %w", err))
	}

	mock := &MockDynamo{
		GetItemOutput: getItemOutput,
		QueryOutput:   queryOutput,
		emulatingErrors: map[FailureCondition]error{
			NotFound:   ErrForceNotFound,
			UserExists: ErrExistingUser,
			None:       nil,
		},
	}

	storage.DefaultClient = mock

	return mock
}

func (d *MockDynamo) ActivateForceFailure(condition FailureCondition) {
	d.mockedErr = d.emulatingErrors[condition]
}

func (d *MockDynamo) DeactivateForceFailure() {
	d.mockedErr = nil
}

func (d *MockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.GetItemOutput{}, d.mockedErr
	}

	if d.GetItemOutput == nil {
		return &dynamodb.GetItemOutput{}, ErrMockNotInitialized
	}

	return d.GetItemOutput, nil
}

func (d *MockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.QueryOutput{}, d.mockedErr
	}

	if d.QueryOutput == nil {
		return &dynamodb.QueryOutput{}, ErrMockNotInitialized
	}

	return d.QueryOutput, nil
}

func (d *MockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.PutItemOutput{}, d.mockedErr
	}

	return &dynamodb.PutItemOutput{}, nil
}

// MockGetItemFromSource mocks the response of the Dynamo DB's GetItem operation using source as the returned item.
func (d *MockDynamo) MockGetItemFromSource(source interface{}) error {
	item, err := attributevalue.MarshalMap(source)
	if err != nil {
		return err
	}

	d.GetItemOutput = &dynamodb.GetItemOutput{
		Item: item,
	}

	return nil
}

// MockQueryFromSource mocks the response of the Dynamo DB's Query operation using source as the returned item.
func (d *MockDynamo) MockQueryFromSource(source interface{}) error {
	item, err := attributevalue.MarshalMap(source)
	if err != nil {
		return err
	}

	d.QueryOutput = &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{item},
	}

	return nil
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
