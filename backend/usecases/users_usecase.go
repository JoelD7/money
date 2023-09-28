package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
)

type UserManager interface {
	GetUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type IncomeGetter interface {
	GetIncomeByPeriod(ctx context.Context, username string, periodID string) ([]*models.Income, error)
}

type ExpenseGetter interface {
	GetExpensesByPeriod(ctx context.Context, username string, periodID string) ([]*models.Expense, error)
}

func NewUserGetter(u UserManager, i IncomeGetter, e ExpenseGetter) func(ctx context.Context, username string) (*models.User, error) {
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

		if newCategory.Budget != nil {
			categoryToUpdate.Budget = newCategory.Budget
		}

		if newCategory.Color != nil {
			categoryToUpdate.Color = newCategory.Color
		}

		newCategories = append(newCategories, categoryToUpdate)

		user.Categories = newCategories

		err = u.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		return nil
	}
}
