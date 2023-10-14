package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
)

type ExpenseManager interface {
	CreateExpense(ctx context.Context, expense *models.Expense) error
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error)
}

func NewExpenseCreator(em ExpenseManager, um UserManager) func(ctx context.Context, username string, expense *models.Expense) error {
	return func(ctx context.Context, username string, expense *models.Expense) error {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return err
		}

		expense.ExpenseID = generateDynamoID("EX")
		expense.Username = username
		expense.Period = user.CurrentPeriod

		err = em.CreateExpense(ctx, expense)
		if err != nil {
			return err
		}

		return nil
	}
}

func NewExpenseUpdater(em ExpenseManager) func(ctx context.Context, expenseID, username string, expense *models.Expense) error {
	return func(ctx context.Context, expenseID, username string, expense *models.Expense) error {
		expense.Username = username
		expense.ExpenseID = expenseID

		err := em.UpdateExpense(ctx, expense)
		if err != nil {
			return err
		}

		return nil
	}
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

func NewExpensesGetter(em ExpenseManager, um UserManager) func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error) {
		expenses, nextKey, err := em.GetExpenses(ctx, username, startKey, pageSize)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(ctx, username, um, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func NewExpensesByCategoriesGetter(em ExpenseManager, um UserManager) func(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
		expenses, nextKey, err := em.GetExpensesByCategory(ctx, username, startKey, categories, pageSize)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(ctx, username, um, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func NewExpensesByPeriodGetter(em ExpenseManager, um UserManager) func(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error) {
		expenses, nextKey, err := em.GetExpensesByPeriod(ctx, username, periodID, startKey, pageSize)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(ctx, username, um, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func NewExpensesByPeriodAndCategoriesGetter(em ExpenseManager, um UserManager) func(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
		expenses, nextKey, err := em.GetExpensesByPeriodAndCategories(ctx, username, periodID, startKey, categories, pageSize)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(ctx, username, um, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func setExpensesCategoryNames(ctx context.Context, username string, um UserManager, expenses []*models.Expense) error {
	user, err := um.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
	}

	categoryNamesByID := make(map[string]string)

	for _, category := range user.Categories {
		categoryNamesByID[category.ID] = *category.Name
	}

	for _, expense := range expenses {
		expense.CategoryName = categoryNamesByID[*expense.CategoryID]
	}

	return nil
}
