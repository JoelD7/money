package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
)

type ExpenseManager interface {
	GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error)
}

func NewExpenseGetter(em ExpenseManager, um UserManager) func(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	return func(ctx context.Context, username, expenseID string) (*models.Expense, error) {
		expense, err := em.GetExpense(ctx, username, expenseID)
		if err != nil {
			return nil, err
		}

		err = setExpensesCategoryNames(ctx, username, um, []*models.Expense{expense})
		if err != nil {
			return expense, err
		}

		return expense, nil
	}
}

func setExpensesCategoryNames(ctx context.Context, username string, um UserManager, expenses []*models.Expense) error {
	user, err := um.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("get user failed for setting category names failed: %v", err)
	}

	categoryNamesByID := make(map[string]string)

	for _, category := range user.Categories {
		categoryNamesByID[category.ID] = *category.Name
	}

	for _, expense := range expenses {
		expense.CategoryName = categoryNamesByID[expense.CategoryID]
	}

	return nil
}
