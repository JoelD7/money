package savingoal

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
}
