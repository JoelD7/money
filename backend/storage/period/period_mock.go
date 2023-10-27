package period

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

var (
	startDate     = models.ToPeriodTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	endDate       = models.ToPeriodTime(time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC))
	defaultPeriod = &models.Period{
		ID:        "2020-01",
		Username:  "test@gmail.com",
		Name:      "January 2020",
		StartDate: startDate,
		EndDate:   endDate,
	}
)

type DynamoMock struct {
	mockedErr      error
	mockedExpenses []*models.Expense
}

func NewDynamoMock() *DynamoMock {
	return &DynamoMock{
		mockedErr: nil,
	}
}

func (d *DynamoMock) ActivateForceFailure(err error) {
	d.mockedErr = err
}

// DeactivateForceFailure deactivates the failures of Dynamo operations.
func (d *DynamoMock) DeactivateForceFailure() {
	d.mockedErr = nil
}

func (d *DynamoMock) CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return period, nil
}

func (d *DynamoMock) GetPeriod(ctx context.Context, username, period string) (*models.Period, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return defaultPeriod, nil
}

func (d *DynamoMock) GetLastPeriod(ctx context.Context, username string) (*models.Period, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return defaultPeriod, nil
}

func (d *DynamoMock) GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, error) {
	if d.mockedErr != nil {
		return nil, d.mockedErr
	}

	return []*models.Period{defaultPeriod}, nil
}
