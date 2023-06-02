package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/JoelD7/money/backend/storage/invalidtoken"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	client *redis.Client

	ErrTokensNotFound = errors.New("no invalid tokens found")
	ErrInvalidTTL     = errors.New("TTL is from a past datetime")
	ErrNoSuchKey      = errors.New("key does not exist")
)

const (
	email     = "test@gmail.com"
	redisURL  = "redis://default:810cc997ccd745debfbbdb567631a5c2@us1-polished-shrew-39844.upstash.io:39844"
	keyPrefix = "invalid_tokens:"
)

func init() {

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt)
}

func GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error) {
	key := keyPrefix + email

	dataStr, err := client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrTokensNotFound
	}

	invalidTokens := make([]*models.InvalidToken, 0)

	err = json.Unmarshal([]byte(dataStr), &invalidTokens)
	if err != nil {
		return nil, err
	}

	if len(invalidTokens) == 0 {
		return nil, ErrTokensNotFound
	}

	return invalidTokens, nil
}

func AddInvalidToken(ctx context.Context, email, token string, ttl int64) error {
	if time.Now().Unix() > ttl {
		return ErrInvalidTTL
	}

	key := keyPrefix + email

	invalidTokens, err := GetInvalidTokens(ctx, email)
	if err != nil {
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

	_, err = client.Set(ctx, key, result, 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, key string) error {
	result, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		return ErrNoSuchKey
	}

	return nil
}

func migrateTokens() {
	ctx := context.Background()

	invalidTokens, err := invalidtoken.GetAllForPerson(ctx, email)
	if err != nil {
		panic(err)
	}

	result, err := utils.GetJsonString(invalidTokens)
	if err != nil {
		panic(err)
	}

	_, err = client.SetNX(ctx, "invalid_tokens:test@gmail.com", result, 0).Result()
	if err != nil {
		panic(err)
	}
}
