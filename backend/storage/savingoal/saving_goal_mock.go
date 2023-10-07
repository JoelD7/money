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

	return &models.SavingGoal{
		Username:     username,
		SavingGoalID: savingGoalID,
		Name:         "mocked_name",
		Goal:         1500,
		Deadline:     time.Now().Add(time.Hour * 24 * 30 * 6),
	}, nil
}

func (m *Mock) GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error) {
	if m.mockedErr != nil {
		return nil, m.mockedErr
	}

	return []*models.SavingGoal{
		{
			Username:     username,
			SavingGoalID: "savingGoalID",
			Name:         "mocked_name",
			Goal:         1500,
			Deadline:     time.Now().Add(time.Hour * 24 * 30 * 6),
		},
	}, nil
}
