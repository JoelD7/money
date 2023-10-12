package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type ExpenseManager interface {
	GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error)
}
