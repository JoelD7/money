package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	client *redis.Client

	redisURL = env.GetString("REDIS_URL", "redis://default:810cc997ccd745debfbbdb567631a5c2@us1-polished-shrew-39844.upstash.io:4000")

	ErrNotFound = errors.New("cache: key not found")
)

const (
	retries       = 3
	backoffFactor = 2
)

type RedisClient struct {
	client *redis.Client
}

func init() {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt)
}

func NewClient() *RedisClient {
	return &RedisClient{client}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}

	backoff := time.Second * 2

	for i := 0; i < retries && err != nil; i++ {
		time.Sleep(backoff)

		fmt.Println(fmt.Sprintf("Retry %d, backoff %s", i+1, backoff.String()))
		value, err = rc.client.Get(ctx, key).Result()
		backoff *= backoffFactor
	}

	return value, err
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_, err := rc.client.Set(ctx, key, value, expiration).Result()

	backoff := time.Second * 2

	for i := 0; i < retries && err != nil; i++ {
		time.Sleep(backoff)

		_, err = rc.client.Set(ctx, key, value, expiration).Result()
		backoff *= backoffFactor
	}

	return err
}
