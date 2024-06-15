package expenses_recurring

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type ExpensesRecurringRepository interface {
	CreateExpenseRecurring(ctx context.Context, expenseRecurring *models.ExpenseRecurring) (*models.ExpenseRecurring, error)
	ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error)
}
