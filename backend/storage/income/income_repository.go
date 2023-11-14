package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)
	GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error)
	GetIncomeByPeriod(ctx context.Context, username, periodID string) ([]*models.Income, error)
}
