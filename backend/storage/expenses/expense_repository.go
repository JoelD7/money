package expenses

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type Repository interface {
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
	BatchDeleteExpenses(ctx context.Context, expenses []*models.Expense) error
}
