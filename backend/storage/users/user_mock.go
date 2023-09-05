package users

import (
	"context"
	"github.com/JoelD7/money/backend/models"
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

type DynamoMock struct {
	mockedErr  error
	mockedUser *models.User
}

func NewDynamoMock() *DynamoMock {
	return &DynamoMock{
		mockedUser: GetDummyUser(),
		mockedErr:  nil,
	}
}

func (d *DynamoMock) SetMockedUser(user *models.User) {
	d.mockedUser = user
}

func (d *DynamoMock) CreateUser(ctx context.Context, fullName, email, password string) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
}

func (d *DynamoMock) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedUser == nil {
		return nil, models.ErrUserNotFound
	}

	return d.mockedUser, nil
}

func (d *DynamoMock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedUser == nil {
		return nil, models.ErrUserNotFound
	}

	return d.mockedUser, nil
}

func (d *DynamoMock) UpdateUser(ctx context.Context, user *models.User) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
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

// GetDummyUser returns the mock item for the user table
func GetDummyUser() *models.User {
	return &models.User{
		UserID:        "123",
		FullName:      "Joel",
		Email:         "test@gmail.com",
		CurrentPeriod: "2023-5",
		Password:      "$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe",
		AccessToken:   hashedDummyToken,
		RefreshToken:  hashedDummyToken,
	}
}
