package users

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	DummyToken         = "header.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjo5OTk5OTk5OTk5fQ.signature"
	DummyPreviousToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0QGdtYWlsLmNvbSIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.nL-Ir6ZnsMHZa7YYwjfpy1QJ1OTmBCFHCDXVSToXUqf3DHA5oWnBtlBuUZ1xTHa5ArQf5vQQIOIrW6p6OjtMdHO3h3-TWOWJIJhbkEmUjS5EMRtZfLWnf9gDnF7CxmUn0yA1qK0B4Nqx57lsI8eMeZKDvN8bqfwlEe53Qy8tYXP5jNxP2zA6Mt7ROCGrfvulTyM0ZwV7klArEKs485NPao8BlyV90s-whjk6h1_mtderbMA2iRxkoARzPRnSftULDYmzCJ3i4IOX9p6xyOcgwecpn93-ya1x1nZtoITZ2It5SYUcrsQ2KhiP2c95bFpJTr6A2UcuAz1Y0GguSR2wlw"
)

var (
	// This is the hashed version of the DummyToken variable with the same hash function we use to store the tokens on
	// the DB. We need this variable for the mock because all tokens are stored hashed on the DB.
	hashedDummyToken = "4f7c5d5d43a3c7e28ea09bc73679378151a3e086ad4360e5469423197a62b665"
)

var mockedPerson *models.User

type MockDynamo struct {
	GetItemOutput *dynamodb.GetItemOutput
	QueryOutput   *dynamodb.QueryOutput
	PutItemOutput *dynamodb.PutItemOutput

	mockedErr error
}

func InitDynamoMock() *MockDynamo {
	mockedPerson = GetMockedUser()

	item, err := attributevalue.MarshalMap(mockedPerson)
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

// GetMockedUser returns the mock item for the user table
func GetMockedUser() *models.User {
	return &models.User{
		FullName:     "Joel",
		Email:        "test@gmail.com",
		Password:     "$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe",
		AccessToken:  hashedDummyToken,
		RefreshToken: hashedDummyToken,
	}
}
