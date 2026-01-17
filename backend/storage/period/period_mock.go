package period

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

var (
	name          = "January 2020"
	defaultPeriod = &models.Period{
		ID:        "2020-01",
		Username:  "test@gmail.com",
		Name:      &name,
		StartDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
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

func (d *DynamoMock) UpdatePeriod(ctx context.Context, period *models.Period) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
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

func (d *DynamoMock) GetPeriods(ctx context.Context, username string, params *models.PeriodQueryParameters) ([]*models.Period, string, error) {
	if d.mockedErr != nil {
		return nil, "", d.mockedErr
	}

	return []*models.Period{defaultPeriod}, "", nil
}

func (d *DynamoMock) BatchGetPeriods(ctx context.Context, username string, periods []string) ([]*models.Period, error) {
	return make([]*models.Period, 0), nil
}

func (d *DynamoMock) DeletePeriod(ctx context.Context, periodID, username string) error {
	if d.mockedErr != nil {
		return d.mockedErr
	}

	return nil
}

func (d *DynamoMock) BatchDeletePeriods(ctx context.Context, periods []*models.Period) error {
	//TODO implement me
	return nil
}

func (d *DynamoMock) GetDefaultPeriod() *models.Period {
	return defaultPeriod
}
