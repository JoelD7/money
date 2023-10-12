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
	GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
	GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error)
}
