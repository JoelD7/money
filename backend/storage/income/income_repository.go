package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type RepositoryAPI interface {
	getIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error)
}

type Repository struct {
	client RepositoryAPI
}

func NewRepository(client RepositoryAPI) *Repository {
	return &Repository{client}
}

func (r *Repository) GetIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error) {
	return r.client.getIncomeByPeriod(ctx, userID, periodID)
}
