package savings

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetSavings(ctx context.Context, username, startKey string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsByPeriod(ctx context.Context, username, startKey, period string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsBySavingGoal(ctx context.Context, startKey, savingGoalID string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsBySavingGoalAndPeriod(ctx context.Context, startKey, savingGoalID, period string, pageSize int) ([]*models.Saving, string, error)
	CreateSaving(ctx context.Context, saving *models.Saving) error
	UpdateSaving(ctx context.Context, saving *models.Saving) error
	DeleteSaving(ctx context.Context, savingID, username string) error
}
