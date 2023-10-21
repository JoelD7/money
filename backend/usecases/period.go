package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type PeriodManager interface {
	CreatePeriod(ctx context.Context, period *models.Period) error
	GetPeriod(ctx context.Context, username, period string) (*models.Period, error)
	GetPeriods(ctx context.Context, username string) ([]*models.Period, error)
}
