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

func (m *Mock) GetSaving(ctx context.Context, username, savingID string) (*models.Saving, error) {
	if m.mockedErr != nil {
		return nil, m.mockedErr
	}

	return m.mockedSavings[0], nil
}

func (m *Mock) GetSavingsByPeriod(ctx context.Context, username, startKey, period string, pageSize int) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	return m.mockedSavings, "next_key", nil
}

func (m *Mock) GetSavingsBySavingGoal(ctx context.Context, startKey, savingGoalID string, pageSize int) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	return m.mockedSavings, "next_key", nil
}

func (m *Mock) GetSavingsBySavingGoalAndPeriod(ctx context.Context, startKey, savingGoalID, period string, pageSize int) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	return m.mockedSavings, "next_key", nil
}

func (m *Mock) GetSavings(ctx context.Context, username, startKey string, pageSize int) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	return m.mockedSavings, "next_key", nil
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
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       250,
		},
		{
			SavingID:     "SV456",
			SavingGoalID: "SVG46",
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       450,
		},
		{
			SavingID:     "SV789",
			SavingGoalID: "SVG789",
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       789,
		},
		{
			SavingID:     "SV159",
			SavingGoalID: "SVG159",
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       156,
		},
	}
}
