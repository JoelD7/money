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

func (d *DynamoMock) getIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedIncome == nil {
		return nil, ErrNotFound
	}

	return d.mockedIncome, nil
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
			UserID:   "test@gmail.com",
			IncomeID: "INC123",
			Amount:   8700,
			Name:     "Salary",
			Date:     time.Date(2023, 5, 15, 20, 0, 0, 0, time.UTC),
			Period:   "2023-5",
		},
		{
			UserID:   "test@gmail.com",
			IncomeID: "INC12",
			Amount:   1500,
			Name:     "Debt collection",
			Date:     time.Date(2023, 5, 15, 20, 0, 0, 0, time.UTC),
			Period:   "2023-5",
		},
	}
}
