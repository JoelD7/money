package period

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error)
	GetPeriod(ctx context.Context, username, period string) (*models.Period, error)
	GetLastPeriod(ctx context.Context, username string) (*models.Period, error)
	GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, error)
}
