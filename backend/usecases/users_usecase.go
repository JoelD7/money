package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
)

type UserGetter interface {
	GetUser(ctx context.Context, username string) (*models.User, error)
}

type IncomeGetter interface {
	GetIncomeByPeriod(ctx context.Context, username string, periodID string) ([]*models.Income, error)
}

type ExpenseGetter interface {
	GetExpensesByPeriod(ctx context.Context, username string, periodID string) ([]*models.Expense, error)
}

func NewUserGetter(u UserGetter, i IncomeGetter, e ExpenseGetter) func(ctx context.Context, username string) (*models.User, error) {
	return func(ctx context.Context, username string) (*models.User, error) {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}

		if user.CurrentPeriod == "" {
			return user, nil
		}

		userExpenses, err := e.GetExpensesByPeriod(ctx, user.Username, user.CurrentPeriod)
		if err != nil {
			return user, fmt.Errorf("the remainder for the user's current period couldn't be calculated: %w", err)
		}

		userIncome, err := i.GetIncomeByPeriod(ctx, user.Username, user.CurrentPeriod)
		if err != nil {
			return user, fmt.Errorf("the remainder for the user's current period couldn't be calculated: %w", err)
		}

		totalExpense := 0.0

		for _, expense := range userExpenses {
			totalExpense += expense.Amount
		}

		totalIncome := 0.0
		for _, inc := range userIncome {
			totalIncome += inc.Amount
		}

		user.Remainder = totalIncome - totalExpense

		return user, nil
	}
}
