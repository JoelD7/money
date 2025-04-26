package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"regexp"
	"sync"
)

var (
	hexColorPattern = "^#[0-9A-Fa-f]{1,6}$"
)

const (
	categoryPrefix = "CTG"
)

type UserManager interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, username string) error
}

type IDGenerator interface {
	GenerateID(prefix string) string
}

func NewUserGetter(u UserManager, i IncomeRepository, e ExpenseManager) func(ctx context.Context, username string) (*models.User, error) {
	return func(ctx context.Context, username string) (*models.User, error) {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}

		if user.CurrentPeriod == "" {
			return user, nil
		}

		userIncome := make([]*models.Income, 0)
		userExpenses := make([]*models.Expense, 0)

		errorCh := make(chan error, 2) // Buffer size is set to the number of goroutines

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			if userIncome, err = getAllIncomeForPeriod(ctx, user.Username, user.CurrentPeriod, i); err != nil {
				errorCh <- err
			}
		}()

		go func() {
			defer wg.Done()
			if userExpenses, err = getAllExpensesForPeriod(ctx, user.Username, user.CurrentPeriod, e); err != nil {
				errorCh <- err
			}
		}()

		wg.Wait()
		close(errorCh)

		for err = range errorCh {
			if err != nil {
				return user, fmt.Errorf("the remainder for the user's current period couldn't be calculated: %w", err)
			}
		}

		totalExpense := 0.0

		for _, expense := range userExpenses {
			totalExpense += *expense.Amount
		}

		totalIncome := 0.0
		for _, inc := range userIncome {
			totalIncome += *inc.Amount
		}

		user.Remainder = totalIncome - totalExpense

		return user, nil
	}
}

func getAllExpensesForPeriod(ctx context.Context, username string, period string, e ExpenseManager) ([]*models.Expense, error) {
	expenses, nextKey, err := e.GetExpensesByPeriod(ctx, username, &models.QueryParameters{Period: period, PageSize: 10})
	if err != nil {
		return nil, fmt.Errorf("get all user expenses failed: %w", err)
	}

	expensesPage := make([]*models.Expense, 0)

	for nextKey != "" {
		expensesPage, nextKey, err = e.GetExpensesByPeriod(ctx, username, &models.QueryParameters{Period: period, PageSize: 10, StartKey: nextKey})
		if err != nil && !errors.Is(err, models.ErrNoMoreItemsToBeRetrieved) {
			return nil, fmt.Errorf("get all user expenses failed: %w", err)
		}

		expenses = append(expenses, expensesPage...)
	}

	return expenses, nil
}

func getAllIncomeForPeriod(ctx context.Context, username string, period string, i IncomeRepository) ([]*models.Income, error) {
	income, nextKey, err := i.GetIncomeByPeriod(ctx, username, &models.QueryParameters{Period: period})
	if err != nil {
		return nil, fmt.Errorf("get all user income failed: %w", err)
	}

	incomePage := make([]*models.Income, 0)

	for nextKey != "" {
		incomePage, nextKey, err = i.GetIncomeByPeriod(ctx, username, &models.QueryParameters{Period: period, StartKey: nextKey})
		if err != nil && !errors.Is(err, models.ErrNoMoreItemsToBeRetrieved) {
			return nil, fmt.Errorf("get all user income failed: %w", err)
		}

		income = append(income, incomePage...)
	}

	return income, nil
}

func NewUserDeleter(u UserManager) func(ctx context.Context, username string) error {
	return func(ctx context.Context, username string) error {
		return u.DeleteUser(ctx, username)
	}
}

func NewCategoryCreator(u UserManager) func(ctx context.Context, username string, category *models.Category) error {
	return func(ctx context.Context, username string, category *models.Category) error {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return err
		}

		if user.Categories == nil {
			user.Categories = make([]*models.Category, 0)
		}

		err = validateCategoryName(category, user.Categories)
		if err != nil {
			return err
		}

		err = validateCategoryColor(category.Color)
		if err != nil {
			return err
		}

		category.ID = generateDynamoID(categoryPrefix)
		user.Categories = append(user.Categories, category)

		return u.UpdateUser(ctx, user)
	}
}

func NewCategoriesGetter(u UserManager) func(ctx context.Context, username string) ([]*models.Category, error) {
	return func(ctx context.Context, username string) ([]*models.Category, error) {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}

		if len(user.Categories) == 0 {
			return nil, models.ErrCategoriesNotFound
		}

		return user.Categories, nil
	}
}

func NewCategoryUpdater(u UserManager) func(ctx context.Context, username, categoryID string, newCategory *models.Category) error {
	return func(ctx context.Context, username, categoryID string, newCategory *models.Category) error {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return err
		}

		err = validateCategoryName(newCategory, user.Categories)
		if err != nil {
			return err
		}

		err = validateCategoryColor(newCategory.Color)
		if err != nil {
			return err
		}

		if user.Categories == nil {
			return models.ErrCategoryNotFound
		}

		newCategories := make([]*models.Category, 0, len(user.Categories))
		var categoryToUpdate *models.Category

		for _, cat := range user.Categories {
			if cat.ID == categoryID {
				categoryToUpdate = cat
				continue
			}

			newCategories = append(newCategories, cat)
		}

		if categoryToUpdate == nil {
			return models.ErrCategoryNotFound
		}

		if newCategory.Name != nil {
			categoryToUpdate.Name = newCategory.Name
		}

		if newCategory.Budget != nil {
			categoryToUpdate.Budget = newCategory.Budget
		}

		if newCategory.Color != nil {
			categoryToUpdate.Color = newCategory.Color
		}

		newCategories = append(newCategories, categoryToUpdate)

		user.Categories = newCategories

		return u.UpdateUser(ctx, user)
	}
}

func validateCategoryName(newCategory *models.Category, userCategories []*models.Category) error {
	for _, category := range userCategories {
		if category.Name != nil && *category.Name == *newCategory.Name {
			return models.ErrCategoryNameAlreadyExists
		}
	}

	return nil
}

func validateCategoryColor(color *string) error {
	if color == nil {
		return nil
	}

	regExp := regexp.MustCompile(hexColorPattern)

	if !regExp.MatchString(*color) {
		return models.ErrInvalidHexColor
	}

	return nil
}
