package cache

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAddInvalidToken(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	email := "test@gmail.com"

	token := "random_token"
	tokenDuration := time.Second * 1

	ttl := time.Now().Add(tokenDuration).Unix()
	err := AddInvalidToken(ctx, email, token, ttl)
	c.Nil(err)

	invalidTokens, err := GetInvalidTokens(ctx, email)
	c.Nil(err)
	c.True(isTokenInvalid(token, invalidTokens))

	// wait for the token to expire
	time.Sleep(time.Second * 2)

	err = AddInvalidToken(ctx, email, "another_token", time.Now().Add(time.Hour*1).Unix())
	c.Nil(err)

	invalidTokens, err = GetInvalidTokens(ctx, email)
	c.Nil(err)
	c.False(isTokenInvalid(token, invalidTokens))
}

func TestAddInvalidTokenFailed(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	email := "test@gmail.com"

	t.Run("Invalid TTL", func(t *testing.T) {
		ttl := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local).Unix()

		err := AddInvalidToken(ctx, email, "token", ttl)
		c.ErrorIs(err, ErrInvalidTTL)
	})

	t.Run("No invalid tokens found", func(t *testing.T) {
		ttl := time.Now().Add(time.Hour * 1).Unix()

		err := AddInvalidToken(ctx, "random", "token", ttl)
		c.ErrorIs(err, ErrTokensNotFound)
	})

	//t.Cleanup(func() {
	//	err := Delete(ctx, keyPrefix+email)
	//	c.Nil(err)
	//})
}

func isTokenInvalid(token string, invalidTokens []*models.InvalidToken) bool {
	for _, it := range invalidTokens {
		if it.Token == token {
			return true
		}
	}

	return false
}
