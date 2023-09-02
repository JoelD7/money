package cache

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

const (
	keyPrefix = "invalid_tokens:"
)

type InvalidTokenCache interface {
	getInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error)
	addInvalidToken(ctx context.Context, email, token string, ttl int64) error
}

type Repository struct {
	client InvalidTokenCache
}

func NewRepository(client InvalidTokenCache) *Repository {
	return &Repository{client}
}

func (r *Repository) GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error) {
	return r.client.getInvalidTokens(ctx, email)
}

func (r *Repository) AddInvalidToken(ctx context.Context, email, token string, ttl int64) error {
	return r.client.addInvalidToken(ctx, email, token, ttl)
}
