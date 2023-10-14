package expenses

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type DynamoMock struct {
	mockedErr      error
	mockedExpenses []*models.Expense
}

func NewDynamoMock() *DynamoMock {
	return &DynamoMock{
		mockedErr:      nil,
		mockedExpenses: GetDummyExpenses(),
	}
}

func (d *DynamoMock) GetExpensesByPeriod(ctx context.Context, username, periodID string) ([]*models.Expense, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, models.ErrExpensesNotFound
	}

	return d.mockedExpenses, nil
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

func (d *DynamoMock) SetMockedExpenses(expenses []*models.Expense) {
	d.mockedExpenses = expenses
}

func GetDummyExpenses() []*models.Expense {
	return []*models.Expense{
		{
			ExpenseID:   "EXP123",
			Username:    "test@mail.com",
			CategoryID:  "",
			Amount:      getFloat64Ptr(893),
			Name:        getStringPtr("Jordan shopping"),
			Notes:       "",
			CreatedDate: time.Date(2023, 5, 12, 20, 15, 0, 0, time.UTC),
			Period:      "2023-5",
			UpdateDate:  time.Time{},
		},
		{
			ExpenseID:   "EXP456",
			Username:    "test@mail.com",
			CategoryID:  "",
			Amount:      getFloat64Ptr(112),
			Name:        getStringPtr("Uber drive"),
			Notes:       "",
			CreatedDate: time.Date(2023, 5, 15, 12, 15, 0, 0, time.UTC),
			Period:      "2023-5",
			UpdateDate:  time.Time{},
		},
		{
			ExpenseID:   "EXP789",
			Username:    "test@mail.com",
			CategoryID:  "",
			Amount:      getFloat64Ptr(525),
			Name:        getStringPtr("Lunch"),
			Notes:       "",
			CreatedDate: time.Date(2023, 5, 12, 11, 15, 0, 0, time.UTC),
			Period:      "2023-5",
			UpdateDate:  time.Time{},
		},
	}
}

func getFloat64Ptr(f float64) *float64 {
	return &f
}

func getStringPtr(s string) *string {
	return &s
}
