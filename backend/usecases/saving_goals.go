package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type SavingGoalManager interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
}

func NewSavingGoalCreator(savingGoalManager SavingGoalManager) func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		return savingGoalManager.CreateSavingGoal(ctx, savingGoal)
	}
}
