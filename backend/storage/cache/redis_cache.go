package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/utils"
)

var (
	ErrInvalidTTL = errors.New("TTL is from a past datetime")
)

type redisCache struct{}

func NewRedisCache() *redisCache {
	return &redisCache{}
}

func (r *redisCache) GetInvalidTokens(ctx context.Context, username string) ([]*models.InvalidToken, error) {
	key := keyPrefix + username

	dataStr, err := redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	invalidTokens := make([]*models.InvalidToken, 0)

	err = json.Unmarshal([]byte(dataStr), &invalidTokens)
	if err != nil {
		return nil, err
	}

	if len(invalidTokens) == 0 {
		return nil, models.ErrInvalidTokensNotFound
	}

	return invalidTokens, nil
}

func (r *redisCache) AddInvalidToken(ctx context.Context, username, token string, ttl int64) error {
	if time.Now().Unix() > ttl && ttl > 0 {
		return ErrInvalidTTL
	}

	key := keyPrefix + username

	invalidTokens, err := r.GetInvalidTokens(ctx, username)
	if err != nil && !errors.Is(err, models.ErrInvalidTokensNotFound) {
		return err
	}

	newInvalidTokens := make([]*models.InvalidToken, 0)
	newInvalidTokens = append(newInvalidTokens, &models.InvalidToken{Token: token, Expire: ttl, CreatedDate: time.Now()})

	now := time.Now().Unix()

	for _, it := range invalidTokens {
		if now <= it.Expire {
			newInvalidTokens = append(newInvalidTokens, it)
		}
	}

	result, err := utils.GetJsonString(newInvalidTokens)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, key, result, 0)
	if err != nil {
		return err
	}

	return nil
}
