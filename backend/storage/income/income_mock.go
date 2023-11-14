package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type DynamoMock struct {
	mockedErr    error
	mockedIncome []*models.Income
}

func NewDynamoMock() *DynamoMock {
	return &DynamoMock{
		mockedErr:    nil,
		mockedIncome: GetDummyIncome(),
	}
}

func (d *DynamoMock) CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return income, nil
}

func (d *DynamoMock) GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return d.mockedIncome[0], nil
}

func (d *DynamoMock) GetIncomeByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Income, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	if d.mockedIncome == nil {
		return nil, "", models.ErrIncomeNotFound
	}

	return d.mockedIncome, "", nil
}

func (d *DynamoMock) SetMockedIncome(income []*models.Income) {
	d.mockedIncome = income
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

func GetDummyIncome() []*models.Income {
	return []*models.Income{
		{
			Username:    "test@gmail.com",
			IncomeID:    "INC123",
			Amount:      getFloatPtr(8700),
			Name:        getStringPtr("Salary"),
			CreatedDate: time.Date(2023, 5, 15, 20, 0, 0, 0, time.UTC),
			Period:      getStringPtr("2023-5"),
		},
		{
			Username:    "test@gmail.com",
			IncomeID:    "INC12",
			Amount:      getFloatPtr(1500),
			Name:        getStringPtr("Debt collection"),
			CreatedDate: time.Date(2023, 5, 15, 20, 0, 0, 0, time.UTC),
			Period:      getStringPtr("2023-5"),
		},
	}
}

func getStringPtr(s string) *string {
	return &s
}

func getFloatPtr(f float64) *float64 {
	return &f
}
