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

type RepositoryAPI interface {
	getExpensesByPeriod(ctx context.Context, userID, periodID string) ([]*models.Expense, error)
}

type Repository struct {
	client RepositoryAPI
}

func NewRepository(client RepositoryAPI) *Repository {
	return &Repository{client}
}

func (r *Repository) GetExpensesByPeriod(ctx context.Context, userID, periodID string) ([]*models.Expense, error) {
	return r.client.getExpensesByPeriod(ctx, userID, periodID)
}
