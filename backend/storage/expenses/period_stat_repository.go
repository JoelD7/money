package expenses

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetPeriodStat(ctx context.Context, period, username, categoryID string) (*models.PeriodStat, error)
}
