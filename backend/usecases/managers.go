package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"net/http"
)

// Expenses

type ExpenseManager interface {
	CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error)
	BatchCreateExpenses(ctx context.Context, expenses []*models.Expense) error

	GetExpenses(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error)
	GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error)
	GetAllExpensesBetweenDates(ctx context.Context, username, startDate, endDate string) ([]*models.Expense, error)
	GetAllExpensesByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, error)

	UpdateExpense(ctx context.Context, expense *models.Expense) error
	BatchUpdateExpenses(ctx context.Context, expenses []*models.Expense) error

	DeleteExpense(ctx context.Context, expenseID, username string) error
}

type ExpenseRecurringManager interface {
	DeleteExpenseRecurring(ctx context.Context, expenseRecurringID, username string) error
}

// User

type UserManager interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, username string) error
}

type InvalidTokenCache interface {
	GetInvalidTokens(ctx context.Context, username string) ([]*models.InvalidToken, error)
	AddInvalidToken(ctx context.Context, username, token string, ttl int64) error
}

type SecretManager interface {
	GetSecret(ctx context.Context, name string) (string, error)
}

// Income

type IncomeRepository interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)

	GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error)
	GetAllIncome(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error)
	GetIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error)
	GetAllIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, error)
	GetAllIncomePeriods(ctx context.Context, username string) ([]string, error)
}

type IncomePeriodCacheManager interface {
	AddIncomePeriods(ctx context.Context, username string, periods []string) error
	GetIncomePeriods(ctx context.Context, username string) ([]string, error)
	DeleteIncomePeriods(ctx context.Context, username string, periods ...string) error
}

type ResourceCacheManager interface {
	AddResource(ctx context.Context, key string, resource interface{}, ttl int64) error
	GetResource(ctx context.Context, key string) (string, error)
}

type JWKSGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type PeriodManager interface {
	CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error)
	UpdatePeriod(ctx context.Context, period *models.Period) error
	GetPeriod(ctx context.Context, username, period string) (*models.Period, error)
	GetLastPeriod(ctx context.Context, username string) (*models.Period, error)
	GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error)
	DeletePeriod(ctx context.Context, periodID, username string) error
}

type SavingGoalManager interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	UpdateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error)
	DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error
	GetAllRecurringSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
}

type SavingsManager interface {
	CreateSaving(ctx context.Context, saving *models.Saving) (*models.Saving, error)
	BatchCreateSavings(ctx context.Context, savings []*models.Saving) error

	GetSaving(ctx context.Context, username, savingID string) (*models.Saving, error)
	GetSavings(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsBySavingGoal(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error)
	GetSavingsBySavingGoalAndPeriod(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error)

	UpdateSaving(ctx context.Context, saving *models.Saving) error

	DeleteSaving(ctx context.Context, savingID, username string) error
}
