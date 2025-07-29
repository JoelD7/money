package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"math"
	"math/rand"
	"time"
)

func NewExpenseCreator(em ExpenseManager, pm PeriodManager, cache ResourceCacheManager) func(ctx context.Context, username, idempotencyKey string, expense *models.Expense) (*models.Expense, error) {
	return func(ctx context.Context, username, idempotencyKey string, expense *models.Expense) (*models.Expense, error) {
		return CreateResource(ctx, cache, idempotencyKey, func() (*models.Expense, error) {
			err := validateExpensePeriod(ctx, expense, username, pm)
			if err != nil {
				return nil, err
			}

			expense.ExpenseID = generateDynamoID("EX")
			expense.Username = username
			expense.CreatedDate = time.Now()

			newExpense, err := em.CreateExpense(ctx, expense)
			if err != nil {
				return nil, err
			}

			return newExpense, nil
		})
	}
}

// Returns an error on 50% of the calls
func randomErr[R models.Resource](resource R) (R, error) {
	if rand.Intn(10)%2 == 0 {
		return resource, nil
	}

	return nil, fmt.Errorf("random error")
}

func NewBatchExpensesCreator(em ExpenseManager) func(ctx context.Context, expenses []*models.Expense) error {
	return func(ctx context.Context, expenses []*models.Expense) error {
		for _, expense := range expenses {
			expense.ExpenseID = generateDynamoID("EX")
			expense.CreatedDate = time.Now()
		}

		return em.BatchCreateExpenses(ctx, expenses)
	}
}

func NewExpenseUpdater(em ExpenseManager, pm PeriodManager, um UserManager) func(ctx context.Context, expenseID, username string, expense *models.Expense) (*models.Expense, error) {
	return func(ctx context.Context, expenseID, username string, expense *models.Expense) (*models.Expense, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
		}

		expense.Username = username
		expense.ExpenseID = expenseID
		expense.UpdateDate = time.Now()

		err = validateExpensePeriod(ctx, expense, username, pm)
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

		err = setExpensesCategoryNames(user, []*models.Expense{updatedExpense})
		if err != nil {
			return updatedExpense, err
		}

		return updatedExpense, nil
	}
}

func NewExpenseGetter(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	return func(ctx context.Context, username, expenseID string) (*models.Expense, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
		}

		expense, err := em.GetExpense(ctx, username, expenseID)
		if err != nil {
			return nil, err
		}

		err = setExpensesPeriods(ctx, username, pm, expense)
		if err != nil {
			return nil, err
		}

		err = setExpensesCategoryNames(user, []*models.Expense{expense})
		if err != nil {
			return expense, err
		}

		return expense, nil
	}
}

func NewExpensesGetter(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
		}

		expenses, nextKey, err := em.GetExpenses(ctx, username, params)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesPeriods(ctx, username, pm, expenses...)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(user, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func setExpensesPeriods(ctx context.Context, username string, pm PeriodManager, expenses ...*models.Expense) error {
	expensesPeriods := make([]string, 0, len(expenses))
	for _, expense := range expenses {
		expensesPeriods = append(expensesPeriods, expense.PeriodID)
	}

	periodResults, err := pm.BatchGetPeriods(ctx, username, expensesPeriods)
	if err != nil {
		return err
	}

	periodsByID := make(map[string]*models.Period, len(periodResults))
	for _, period := range periodResults {
		periodsByID[period.ID] = period
	}

	for _, expense := range expenses {
		expense.PeriodName = *periodsByID[expense.PeriodID].Name
	}

	return nil
}

func NewExpensesByCategoriesGetter(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
		}

		expenses, nextKey, err := em.GetExpensesByCategory(ctx, username, params)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesPeriods(ctx, username, pm, expenses...)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(user, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func NewExpensesByPeriodGetter(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, "", fmt.Errorf("%v", err)
		}

		expenses, nextKey, err := em.GetExpensesByPeriod(ctx, username, params)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesPeriods(ctx, username, pm, expenses...)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(user, expenses)
		if err != nil {
			return expenses, "", err
		}

		return expenses, nextKey, nil
	}
}

func NewExpensesByPeriodAndCategoriesGetter(em ExpenseManager, um UserManager, pm PeriodManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
		user, err := um.GetUser(ctx, username)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", models.ErrCategoryNameSettingFailed, err)
		}

		expenses, nextKey, err := em.GetExpensesByPeriodAndCategories(ctx, username, params)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesPeriods(ctx, username, pm, expenses...)
		if err != nil {
			return nil, "", err
		}

		err = setExpensesCategoryNames(user, expenses)
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

func NewExpensesPeriodSetter(em ExpenseManager, pm PeriodManager) func(ctx context.Context, username, periodID string) error {
	return func(ctx context.Context, username, periodID string) error {
		period, err := pm.GetPeriod(ctx, username, periodID)
		if err != nil {
			return err
		}

		startDate := period.StartDate.Format(time.DateOnly)
		endDate := period.EndDate.Format(time.DateOnly)

		expenses, err := em.GetAllExpensesBetweenDates(ctx, username, startDate, endDate)
		if err != nil {
			return err
		}

		for _, expense := range expenses {
			if expense.PeriodID == "" {
				expense.PeriodID = period.ID
			}
		}

		err = em.BatchUpdateExpenses(ctx, expenses)
		if err != nil {
			return err
		}

		return nil
	}
}

func setExpensesCategoryNames(user *models.User, expenses []*models.Expense) error {
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
	if expense.PeriodID == "" {
		return nil
	}

	var err error

	_, err = p.GetPeriod(ctx, username, expense.PeriodID)
	if errors.Is(err, models.ErrPeriodNotFound) {
		return models.ErrInvalidPeriod
	}

	if err != nil {
		return fmt.Errorf("check if expense period is valid failed: %v", err)
	}

	return nil
}

func NewExpenseRecurringEliminator(em ExpenseRecurringManager) func(ctx context.Context, expenseRecurringID, username string) error {
	return func(ctx context.Context, expenseRecurringID, username string) error {
		return em.DeleteExpenseRecurring(ctx, expenseRecurringID, username)
	}
}

func NewCategoryExpenseSummaryGetter(em ExpenseManager) func(ctx context.Context, username, periodID string) ([]*models.CategoryExpenseSummary, error) {
	return func(ctx context.Context, username, periodID string) ([]*models.CategoryExpenseSummary, error) {
		expenses, err := em.GetAllExpensesByPeriod(ctx, username, &models.QueryParameters{Period: periodID})
		if err != nil {
			return nil, err
		}

		categoryExpenses := make([]*models.CategoryExpenseSummary, 0)
		totalExpensesByCategory := make(map[string]float64)

		for _, expense := range expenses {
			if expense.CategoryID != nil && expense.Amount != nil {
				totalExpensesByCategory[*expense.CategoryID] += *expense.Amount
			}
		}

		for categoryID, total := range totalExpensesByCategory {
			categoryExpenses = append(categoryExpenses, &models.CategoryExpenseSummary{
				CategoryID: categoryID,
				Total:      math.Round(total*100) / 100,
				Period:     periodID,
			})
		}

		return categoryExpenses, nil
	}
}
