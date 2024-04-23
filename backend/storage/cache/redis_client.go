package cache

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	redisClient RedisAPI

	redisURL = env.GetString("REDIS_URL", "redis://default:d45d4e1dbf8a4b809be74bda839a7a80@us1-casual-teal-41561.upstash.io:41561")
)

type RedisAPI interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type RedisClient struct {
	client *redis.Client
}

func init() {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	opt.ContextTimeoutEnabled = true

	redisClient = &RedisClient{redis.NewClient(opt)}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", models.ErrInvalidTokensNotFound
	}

	return value, err
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := rc.client.Set(ctx, key, value, expiration).Result()

	return err
}
