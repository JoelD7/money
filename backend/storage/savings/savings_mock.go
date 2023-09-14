package savings

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"strings"
	"time"
)

type Mock struct {
	mockedErr     error
	mockedSavings []*models.Saving
}

func NewMock() *Mock {
	return &Mock{
		mockedSavings: GetDummySavings(),
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

func (m *Mock) CreateSaving(ctx context.Context, saving *models.Saving) error {
	if m.mockedErr != nil {
		return m.mockedErr
	}

	m.mockedSavings = append(m.mockedSavings, saving)

	return nil
}

func (m *Mock) UpdateSaving(ctx context.Context, saving *models.Saving) error {
	if m.mockedErr != nil && strings.Contains(m.mockedErr.Error(), "ConditionalCheckFailedException") {
		return models.ErrUpdateSavingNotFound
	}

	if m.mockedErr != nil {
		return m.mockedErr
	}

	return nil
}

func (m *Mock) DeleteSaving(ctx context.Context, savingID, email string) error {
	if m.mockedErr != nil && strings.Contains(m.mockedErr.Error(), "ConditionalCheckFailedException") {
		return models.ErrDeleteSavingNotFound
	}

	if m.mockedErr != nil {
		return m.mockedErr
	}

	return nil
}

func GetDummySavings() []*models.Saving {
	return []*models.Saving{
		{
			SavingID:     "SV123",
			SavingGoalID: "SVG123",
			Email:        "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       250,
		},
		{
			SavingID:     "SV456",
			SavingGoalID: "SVG46",
			Email:        "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       450,
		},
		{
			SavingID:     "SV789",
			SavingGoalID: "SVG789",
			Email:        "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       789,
		},
		{
			SavingID:     "SV159",
			SavingGoalID: "SVG159",
			Email:        "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       156,
		},
	}
}
