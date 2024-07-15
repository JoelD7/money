package expenses

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
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

func (d *DynamoMock) CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return expense, nil
}

func (d *DynamoMock) BatchCreateExpenses(ctx context.Context, log logger.LogAPI, expenses []*models.Expense) error {
	//TODO implement me
	return nil
}

func (d *DynamoMock) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
}

func (d *DynamoMock) GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, "", models.ErrExpensesNotFound
	}

	return d.mockedExpenses, "", nil
}

func (d *DynamoMock) GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, "", models.ErrExpensesNotFound
	}

	return d.mockedExpenses, "", nil
}

func (d *DynamoMock) GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, "", models.ErrExpensesNotFound
	}

	return d.mockedExpenses, "", nil
}

func (d *DynamoMock) GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, "", models.ErrExpensesNotFound
	}

	return d.mockedExpenses, "", nil
}

func (d *DynamoMock) GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	if d.mockedExpenses == nil {
		return nil, models.ErrExpenseNotFound
	}

	return d.mockedExpenses[0], nil
}

func (d *DynamoMock) DeleteExpense(ctx context.Context, expenseID, username string) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
}

func (d *DynamoMock) GetAllExpensesBetweenDates(ctx context.Context, username, startDate, endDate string) ([]*models.Expense, error) {
	//TODO implement me
	return nil, nil
}

func (d *DynamoMock) BatchUpdateExpenses(ctx context.Context, log logger.LogAPI, expenses []*models.Expense) error {
	//TODO implement me
	return nil
}

func (d *DynamoMock) BatchDeleteExpenses(ctx context.Context, expenses []*models.Expense) error {
	//TODO implement me
	return nil
}

func GetDummyExpenses() []*models.Expense {
	return []*models.Expense{
		{
			ExpenseID:   "EXP123",
			Username:    "test@mail.com",
			CategoryID:  getStringPtr(""),
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
			CategoryID:  getStringPtr(""),
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
			CategoryID:  getStringPtr(""),
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
