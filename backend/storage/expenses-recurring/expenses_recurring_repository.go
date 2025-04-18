package expenses_recurring

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateExpenseRecurring(ctx context.Context, expenseRecurring *models.ExpenseRecurring) (*models.ExpenseRecurring, error)
	BatchCreateExpenseRecurring(ctx context.Context, expenseRecurring []*models.ExpenseRecurring) error

	ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error)
	GetExpenseRecurring(ctx context.Context, expenseRecurringID, username string) (*models.ExpenseRecurring, error)

	BatchDeleteExpenseRecurring(ctx context.Context, expenseRecurring []*models.ExpenseRecurring) error
	DeleteExpenseRecurring(ctx context.Context, expenseRecurringID, username string) error
}
