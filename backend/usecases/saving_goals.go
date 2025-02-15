package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type SavingGoalManager interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
	DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error
}

func NewSavingGoalCreator(savingGoalManager SavingGoalManager) func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		return savingGoalManager.CreateSavingGoal(ctx, savingGoal)
	}
}

func NewSavingGoalGetter(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
	return func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
		return savingGoalManager.GetSavingGoal(ctx, username, savingGoalID)
	}
}

func NewSavingGoalEliminator(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string) error {
	return func(ctx context.Context, username, savingGoalID string) error {
		return savingGoalManager.DeleteSavingGoal(ctx, username, savingGoalID)
	}
}
