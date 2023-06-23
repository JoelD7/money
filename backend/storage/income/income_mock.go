package income

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type DynamoMock struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput
	PutItemOutput *dynamodb.PutItemOutput

	mockedErr error
}

func InitDynamoMock() *DynamoMock {
	items, err := GetMockedIncomeAsItems()
	if err != nil {
		panic(fmt.Errorf("invalid_token Dynamo mock cannot be initialized: %v", err))
	}

	mock := &DynamoMock{
		GetItemOutput: &dynamodb.GetItemOutput{
			Item: items[0],
		},
		QueryOutput: &dynamodb.QueryOutput{
			Items: items,
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
func (d *DynamoMock) ActivateForceFailure(err error) {
	d.mockedErr = err
}

// DeactivateForceFailure deactivates the failures of Dynamo operations.
func (d *DynamoMock) DeactivateForceFailure() {
	d.mockedErr = nil
}

func (d *DynamoMock) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.GetItemOutput{}, d.mockedErr
	}

	return d.GetItemOutput, nil
}

func (d *DynamoMock) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.QueryOutput{}, d.mockedErr
	}

	return d.QueryOutput, nil
}

func (d *DynamoMock) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if d.mockedErr != nil {
		return &dynamodb.PutItemOutput{}, d.mockedErr
	}

	return d.PutItemOutput, nil
}

// MockGetItemFromSource mocks the response of the Dynamo DB's GetItem operation using source as the returned item.
func (d *DynamoMock) MockGetItemFromSource(source interface{}) error {
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
func (d *DynamoMock) MockQueryFromSource(source interface{}) error {
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
func (d *DynamoMock) EmptyTable() {
	d.QueryOutput = &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{},
	}

	d.GetItemOutput = &dynamodb.GetItemOutput{}
}

func GetMockedIncomeAsItems() ([]map[string]types.AttributeValue, error) {
	incomeList := []*models.Income{
		{
			UserID:   "",
			IncomeID: "INC123",
			Amount:   25000,
			Name:     "Salary",
			Date:     time.Date(2023, 5, 15, 20, 0, 0, 0, nil),
			Period:   "2023-5",
		},
		{
			UserID:   "",
			IncomeID: "INC12",
			Amount:   1500,
			Name:     "Debt collection",
			Date:     time.Date(2023, 5, 15, 20, 0, 0, 0, nil),
			Period:   "2023-5",
		},
	}

	items := make([]map[string]types.AttributeValue, 0)

	for _, e := range incomeList {
		item, err := attributevalue.MarshalMap(e)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
