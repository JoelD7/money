package cache

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

const (
	keyPrefix = "invalid_tokens:"
)

type InvalidTokenManager interface {
	GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error)
	AddInvalidToken(ctx context.Context, email, token string, ttl int64) error
}
