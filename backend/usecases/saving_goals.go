package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type SavingGoalManager interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	UpdateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error)
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

func NewSavingGoalsGetter(savingGoalManager SavingGoalManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
		return savingGoalManager.GetSavingGoals(ctx, username, params)
	}
}

func NewSavingGoalUpdator(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username, savingGoalID string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		savingGoal.SavingGoalID = savingGoalID
		return savingGoalManager.UpdateSavingGoal(ctx, savingGoal)
	}
}

func NewSavingGoalEliminator(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string) error {
	return func(ctx context.Context, username, savingGoalID string) error {
		return savingGoalManager.DeleteSavingGoal(ctx, username, savingGoalID)
	}
}
