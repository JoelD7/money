package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error)
}
