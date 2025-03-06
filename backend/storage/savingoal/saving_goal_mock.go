package savingoal

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type Mock struct {
	mockedErr error
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) ActivateForceFailure(err error) {
	m.mockedErr = err
}

func (m *Mock) DeactivateForceFailure() {
	m.mockedErr = nil
}

func (m *Mock) GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
	if m.mockedErr != nil {
		return nil, m.mockedErr
	}

	name := "mocked_name"
	target := float64(1500)
	deadline := time.Now().Add(time.Hour * 24 * 30 * 6)

	return &models.SavingGoal{
		Username:     username,
		SavingGoalID: savingGoalID,
		Name:         &name,
		Target:       &target,
		Deadline:     &deadline,
	}, nil
}

func (m *Mock) CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return nil, nil
}

func (m *Mock) UpdateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return nil, nil
}

func (m *Mock) DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error {
	return nil
}

func (m *Mock) GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
	if m.mockedErr != nil {
		return nil, "", m.mockedErr
	}

	name := "mocked_name"
	target := float64(1500)
	deadline := time.Now().Add(time.Hour * 24 * 30 * 6)

	return []*models.SavingGoal{
		{
			Username:     username,
			SavingGoalID: "savingGoalID",
			Name:         &name,
			Target:       &target,
			Deadline:     &deadline,
		},
	}, "", nil
}
