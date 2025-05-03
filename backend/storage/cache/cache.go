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

type IncomePeriodCacheManager interface {
	AddIncomePeriods(ctx context.Context, username string, periods []string) error
	GetIncomePeriods(ctx context.Context, username string) ([]string, error)
	DeleteIncomePeriods(ctx context.Context, username string, periods ...string) error
}

// IdempotenceCacheManager handles reads and writes to cached resources with idempotency keys
type IdempotenceCacheManager interface {
	// AddResource adds a resource to the cache for ttl seconds
	AddResource(ctx context.Context, key string, resource interface{}, ttl int64) error
	// GetResource gets a resource from the cache
	GetResource(ctx context.Context, key string) (string, error)
}
