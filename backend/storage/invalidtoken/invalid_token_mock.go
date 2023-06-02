package invalidtoken

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type MockDynamo struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput
	PutItemOutput *dynamodb.PutItemOutput

	mockedErr error
}

func InitDynamoMock() *MockDynamo {
	itSample := &models.InvalidToken{
		Token:  "invalidated token",
		Expire: time.Now().Add(time.Hour * 1).Unix(),
	}

	item, err := attributevalue.MarshalMap(itSample)
	if err != nil {
		panic(fmt.Errorf("invalid_token Dynamo mock cannot be initialized: %v", err))
	}

	mock := &MockDynamo{
		GetItemOutput: &dynamodb.GetItemOutput{
			Item: item,
		},
		QueryOutput: &dynamodb.QueryOutput{
			Items: []map[string]types.AttributeValue{item},
		},
		PutItemOutput: &dynamodb.PutItemOutput{},
		mockedErr:     nil,
	}

	DefaultClient = mock

	return mock
}

// ActivateForceFailure makes any of the Dynamo operations fail with the specified error.
// This invocation should always be followed by a deferred call to DeactivateForceFailure so that no other tests are
// affected by this behavior.
func (d *MockDynamo) ActivateForceFailure(err error) {
	d.mockedErr = err
}

// DeactivateForceFailure deactivates the failures of Dynamo operations.
func (d *MockDynamo) DeactivateForceFailure() {
	d.mockedErr = nil
}

func (d *MockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.GetItemOutput{}, d.mockedErr
	}

	return d.GetItemOutput, nil
}

func (d *MockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.QueryOutput{}, d.mockedErr
	}

	return d.QueryOutput, nil
}

func (d *MockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.PutItemOutput{}, d.mockedErr
	}

	return d.PutItemOutput, nil
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

// EmptyTable makes the mocked table to be empty
func (d *MockDynamo) EmptyTable() {
	d.QueryOutput = &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{},
	}

	d.GetItemOutput = &dynamodb.GetItemOutput{}
}
