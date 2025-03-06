package savings

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetSaving(ctx context.Context, username, savingID string) (*models.Saving, error)
	GetSavings(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsBySavingGoal(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsBySavingGoalAndPeriod(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error)
	CreateSaving(ctx context.Context, saving *models.Saving) (*models.Saving, error)
	UpdateSaving(ctx context.Context, saving *models.Saving) error
	DeleteSaving(ctx context.Context, savingID, username string) error
}
