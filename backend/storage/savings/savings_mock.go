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

	for _, saving := range m.mockedSavings {
		if saving.SavingID == savingID {
			return saving, nil
		}
	}

	return nil, models.ErrSavingNotFound
}

func (m *Mock) GetSavingsByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	savings := make([]*models.Saving, 0)

	for _, saving := range m.mockedSavings {
		if *saving.Period == params.Period && saving.Username == username {
			savings = append(savings, saving)
		}
	}

	return savings, "next_key", nil
}

func (m *Mock) GetSavingsBySavingGoal(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	savings := make([]*models.Saving, 0)

	for _, saving := range m.mockedSavings {
		if *saving.SavingGoalID == params.SavingGoalID {
			savings = append(savings, saving)
		}
	}

	return savings, "next_key", nil
}

func (m *Mock) GetSavingsBySavingGoalAndPeriod(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	savings := make([]*models.Saving, 0)
	for _, saving := range m.mockedSavings {
		if *saving.SavingGoalID == params.SavingGoalID && *saving.Period == params.Period {
			savings = append(savings, saving)
		}
	}

	return savings, "next_key", nil
}

func (m *Mock) GetSavings(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	savings := make([]*models.Saving, 0)
	for _, saving := range m.mockedSavings {
		if saving.Username == username {
			savings = append(savings, saving)
		}
	}

	return savings, "next_key", nil
}

func (m *Mock) CreateSaving(ctx context.Context, saving *models.Saving) (*models.Saving, error) {
	if m.mockedErr != nil {
		return nil, m.mockedErr
	}

	m.mockedSavings = append(m.mockedSavings, saving)

	return saving, nil
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

func (m *Mock) BatchCreateSavings(ctx context.Context, savings []*models.Saving) error {
	return nil
}

func GetDummySavings() []*models.Saving {
	return []*models.Saving{
		{
			SavingID:     "SV123",
			SavingGoalID: getStringPtr("SVG123"),
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       getFloatPtr(250),
		},
		{
			SavingID:     "SV456",
			SavingGoalID: getStringPtr("SVG46"),
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       getFloatPtr(450),
		},
		{
			SavingID:     "SV789",
			SavingGoalID: getStringPtr("SVG789"),
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       getFloatPtr(789),
		},
		{
			SavingID:     "SV159",
			SavingGoalID: getStringPtr("SVG159"),
			Username:     "test@gmail.com",
			CreatedDate:  time.Now(),
			Amount:       getFloatPtr(156),
		},
	}
}

func getStringPtr(s string) *string {
	return &s
}

func getFloatPtr(f float64) *float64 {
	return &f
}
