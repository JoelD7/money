package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/utils"
)

var (
	ErrInvalidTTL = errors.New("TTL is from a past datetime")
	redisURL      = env.GetString("REDIS_URL", "redis://default:randome@random.upstash.io:41561")
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	opt.ContextTimeoutEnabled = true

	return &RedisCache{
		redis.NewClient(opt),
	}
}

func (r *RedisCache) GetInvalidTokens(ctx context.Context, username string) ([]*models.InvalidToken, error) {
	key := keyPrefix + username

	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("%w:%v", models.ErrInvalidTokensNotFound, err)
	}

	invalidTokens := make([]*models.InvalidToken, 0)

	err = json.Unmarshal([]byte(value), &invalidTokens)
	if err != nil {
		return nil, err
	}

	if len(invalidTokens) == 0 {
		return nil, models.ErrInvalidTokensNotFound
	}

	return invalidTokens, nil
}

func (r *RedisCache) AddInvalidToken(ctx context.Context, username, token string, ttl int64) error {
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

	_, err = r.client.Set(ctx, key, result, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
