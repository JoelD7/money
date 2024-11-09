package cache

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

const (
	invalidTokenKeyPrefix  = "invalid_tokens"
	incomePeriodsKeyPrefix = "income_periods"
)

type InvalidTokenManager interface {
	GetInvalidTokens(ctx context.Context, username string) ([]*models.InvalidToken, error)
	AddInvalidToken(ctx context.Context, username, token string, ttl int64) error
}
