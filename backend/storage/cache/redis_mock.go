package cache

import (
	"context"
	"errors"
	"time"
)

type RedisMock struct {
	store map[string]string
}

func InitRedisMock() {
	redisClient = &RedisMock{
		store: make(map[string]string),
	}
}

func (r *RedisMock) Get(ctx context.Context, key string) (string, error) {
	value, ok := r.store[key]
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

func (r *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, ok := value.(string)
	if !ok {
		return errors.New("redis mock: parse error")
	}

	r.store[key] = data

	return nil
}
