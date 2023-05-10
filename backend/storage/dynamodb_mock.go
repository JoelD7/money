package storage

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrForceNotFound = errors.New("force not found")
)

type Output struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput
	PutItemOutput *dynamodb.PutItemOutput
}

type MockDynamo struct {
	outputByTable map[string]*Output

	mockedErr error
}

func InitDynamoMock() *MockDynamo {
	mock := &MockDynamo{
		outputByTable: map[string]*Output{},
	}

	DefaultClient = mock

	return mock
}

func (d *MockDynamo) ActivateForceFailure(err error) {
	d.mockedErr = err
}

func (d *MockDynamo) DeactivateForceFailure() {
	d.mockedErr = nil
}

func (d *MockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.GetItemOutput{}, d.mockedErr
	}

	output, ok := d.outputByTable[*params.TableName]
	if !ok {
		return &dynamodb.GetItemOutput{}, nil
	}

	return output.GetItemOutput, nil
}

func (d *MockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.QueryOutput{}, d.mockedErr
	}

	output, ok := d.outputByTable[*params.TableName]
	if !ok {
		return &dynamodb.QueryOutput{}, nil
	}

	return output.QueryOutput, nil
}

func (d *MockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.PutItemOutput{}, d.mockedErr
	}

	output, ok := d.outputByTable[*params.TableName]
	if !ok {
		return &dynamodb.PutItemOutput{}, nil
	}

	return output.PutItemOutput, nil
}

// MockGetItemFromSource mocks the response of the Dynamo DB's GetItem operation using source as the returned item.
func (d *MockDynamo) MockGetItemFromSource(tableName string, source interface{}) error {
	item, err := attributevalue.MarshalMap(source)
	if err != nil {
		return err
	}

	if d.outputByTable[tableName] == nil {
		d.outputByTable[tableName] = &Output{}
	}

	d.outputByTable[tableName].GetItemOutput = &dynamodb.GetItemOutput{
		Item: item,
	}

	return nil
}

// MockQueryFromSource mocks the response of the Dynamo DB's Query operation using source as the returned item.
func (d *MockDynamo) MockQueryFromSource(tableName string, source interface{}) error {
	item, err := attributevalue.MarshalMap(source)
	if err != nil {
		return err
	}

	if d.outputByTable[tableName] == nil {
		d.outputByTable[tableName] = &Output{}
	}

	d.outputByTable[tableName].QueryOutput = &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{item},
	}

	return nil
}
