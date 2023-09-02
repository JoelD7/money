package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
)

type UserGetter interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type IncomeGetter interface {
	GetIncomeByPeriod(ctx context.Context, userID string, periodID string) ([]*models.Income, error)
}

type ExpenseGetter interface {
	GetExpensesByPeriod(ctx context.Context, userID string, periodID string) ([]*models.Expense, error)
}

func NewUserGetter(u UserGetter, i IncomeGetter, e ExpenseGetter) func(ctx context.Context, email string) (*models.User, error) {
	return func(ctx context.Context, email string) (*models.User, error) {
		user, err := u.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		if user.CurrentPeriod == "" {
			return user, nil
		}

		userExpenses, err := e.GetExpensesByPeriod(ctx, user.UserID, user.CurrentPeriod)
		if err != nil {
			return user, fmt.Errorf("the remainder for the user's current period couldn't be calculated: %w", err)
		}

		userIncome, err := i.GetIncomeByPeriod(ctx, user.UserID, user.CurrentPeriod)
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
