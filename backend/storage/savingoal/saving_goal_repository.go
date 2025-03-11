package savingoal

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	UpdateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error)
	GetAllRecurringSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
	DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error
}
