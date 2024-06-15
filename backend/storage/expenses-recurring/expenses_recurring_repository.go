package expenses_recurring

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type ExpensesRecurringRepository interface {
	ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error)
}
