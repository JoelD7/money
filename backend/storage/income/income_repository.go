package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)
	BatchCreateIncome(ctx context.Context, incomes []*models.Income) error

	GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error)
	GetAllIncome(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, error)
	GetIncomeByPeriod(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, error)
	GetAllIncomeByPeriod(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, error)
	GetAllIncomePeriods(ctx context.Context, username string) ([]string, error)

	BatchDeleteIncome(ctx context.Context, income []*models.Income) error
}
