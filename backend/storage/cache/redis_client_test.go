package cache

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	key := "invalid_tokens:test@gmail.com"

	value, err := redisClient.Get(ctx, key)
	c.Empty(value)
	c.Nil(err)

}
