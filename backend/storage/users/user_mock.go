package users

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	mockedErr   error
	mockedUsers []*models.User
}

func NewDynamoMock() *DynamoMock {
	return &DynamoMock{
		mockedUsers: []*models.User{GetDummyUser()},
		mockedErr:   nil,
	}
}

func (d *DynamoMock) CreateUser(ctx context.Context, user *models.User) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	d.mockedUsers = append(d.mockedUsers, user)

	return nil
}

func (d *DynamoMock) GetUser(ctx context.Context, username string) (*models.User, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedUsers == nil {
		return nil, models.ErrUserNotFound
	}

	for _, user := range d.mockedUsers {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, models.ErrUserNotFound
}

func (d *DynamoMock) UpdateUser(ctx context.Context, user *models.User) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	for _, mockedUser := range d.mockedUsers {
		if mockedUser.Username == user.Username {
			mockedUser.FullName = user.FullName
			mockedUser.Username = user.Username
			mockedUser.CurrentPeriod = user.CurrentPeriod
			mockedUser.Password = user.Password
			mockedUser.AccessToken = user.AccessToken
			mockedUser.RefreshToken = user.RefreshToken
			mockedUser.Categories = user.Categories
			mockedUser.UpdatedDate = user.UpdatedDate
			mockedUser.CreatedDate = user.CreatedDate
			mockedUser.Remainder = user.Remainder
			return nil
		}
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
		FullName:      "Joel",
		Username:      "test@gmail.com",
		CurrentPeriod: "2023-5",
		Password:      "$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe",
		AccessToken:   hashedDummyToken,
		RefreshToken:  hashedDummyToken,
		Categories: []*models.Category{
			{
				ID:    "CTGzJeEzCNz6HMTiPKwgPmj",
				Name:  aws.String("Entertainment"),
				Color: aws.String("#ff8733"),
			},
			{
				ID:    "CTGtClGT160UteOl02jIH4F",
				Name:  aws.String("Health"),
				Color: aws.String("#00b85e"),
			},
			{
				ID:    "CTGrR7fO4ndmI0IthJ7Wg8f",
				Name:  aws.String("Utilities"),
				Color: aws.String("#009eb8"),
			},
		},
	}
}
