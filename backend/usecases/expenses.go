package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type ExpenseManager interface {
	CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error)
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error)
	DeleteExpense(ctx context.Context, expenseID, username string) error
}

func NewExpenseCreator(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username string, expense *models.Expense) (*models.Expense, error) {
	return func(ctx context.Context, username string, expense *models.Expense) (*models.Expense, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}

		err = validateExpensePeriod(ctx, expense, username, pm)
		if err != nil {
			return nil, err
		}

		expense.ExpenseID = generateDynamoID("EX")
		expense.Username = username
		expense.Period = &user.CurrentPeriod
		expense.CreatedDate = time.Now()

		newExpense, err := em.CreateExpense(ctx, expense)
		if err != nil {
			return nil, err
		}

		return newExpense, nil
	}
}

func NewExpenseUpdater(em ExpenseManager, pm PeriodManager, um UserManager) func(ctx context.Context, expenseID, username string, expense *models.Expense) (*models.Expense, error) {
	return func(ctx context.Context, expenseID, username string, expense *models.Expense) (*models.Expense, error) {
		expense.Username = username
		expense.ExpenseID = expenseID
		expense.UpdateDate = time.Now()

		err := validateExpensePeriod(ctx, expense, username, pm)
		if err != nil {
			return nil, err
		}

		err = em.UpdateExpense(ctx, expense)
		if err != nil {
			return nil, err
		}

		updatedExpense, err := em.GetExpense(ctx, username, expenseID)
		if err != nil {
			return nil, fmt.Errorf("getting updated expense failed: %w", err)
		}

		err = setExpensesCategoryNames(ctx, username, um, []*models.Expense{updatedExpense})
		if err != nil {
			return updatedExpense, err
		}

		return updatedExpense, nil
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

func NewExpensesDeleter(em ExpenseManager) func(ctx context.Context, expenseID, username string) error {
	return func(ctx context.Context, expenseID, username string) error {
		return em.DeleteExpense(ctx, expenseID, username)
	}
}

func setExpensesCategoryNames(ctx context.Context, username string, um UserManager, expenses []*models.Expense) error {
	user, err := um.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
	}

	categoryNamesByID := make(map[string]string)

	for _, category := range user.Categories {
		if category.Name != nil {
			categoryNamesByID[category.ID] = *category.Name
		}
	}

	for _, expense := range expenses {
		if expense.CategoryID != nil {
			expense.CategoryName = categoryNamesByID[*expense.CategoryID]
		}
	}

	return nil
}

func validateExpensePeriod(ctx context.Context, expense *models.Expense, username string, p PeriodManager) error {
	if expense.Period == nil {
		return nil
	}

	periods := make([]*models.Period, 0)
	curPeriods := make([]*models.Period, 0)
	nextKey := ""
	var err error

	for {
		curPeriods, nextKey, err = p.GetPeriods(ctx, username, nextKey, 0)
		if err != nil {
			return fmt.Errorf("check if expense period is valid failed: %v", err)
		}

		periods = append(periods, curPeriods...)

		if nextKey == "" {
			break
		}
	}

	for _, period := range periods {
		if period.ID == *expense.Period {
			return nil
		}
	}

	return models.ErrInvalidPeriod
}
