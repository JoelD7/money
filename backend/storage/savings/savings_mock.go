package savings

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type Mock struct {
	mockedErr     error
	mockedSavings []*models.Saving
}

func NewMock() *Mock {
	return &Mock{
		mockedSavings: getDummySavings(),
	}
}

func (m *Mock) ActivateForceFailure(err error) {
	m.mockedErr = err
}

func (m *Mock) DeactivateForceFailure() {
	m.mockedErr = nil
}

func (m *Mock) GetSavings(ctx context.Context, email string) ([]*models.Saving, error) {
	if m.mockedErr != nil {
		return nil, m.mockedErr
	}

	return m.mockedSavings, nil
}

func getDummySavings() []*models.Saving {
	return []*models.Saving{
		{
			SavingID:     "SV123",
			SavingGoalID: "SVG123",
			Email:        "test@gmail.com",
			CreationDate: time.Time{},
			Amount:       250,
		},
		{
			SavingID:     "SV456",
			SavingGoalID: "SVG46",
			Email:        "test@gmail.com",
			CreationDate: time.Time{},
			Amount:       450,
		},
		{
			SavingID:     "SV789",
			SavingGoalID: "SVG789",
			Email:        "test@gmail.com",
			CreationDate: time.Time{},
			Amount:       789,
		},
		{
			SavingID:     "SV159",
			SavingGoalID: "SVG159",
			Email:        "test@gmail.com",
			CreationDate: time.Time{},
			Amount:       156,
		},
	}
}
