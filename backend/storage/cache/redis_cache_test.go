package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildKey(t *testing.T) {
	c := require.New(t)

	key := buildKey("test")
	c.Equal("test:", key)

	key = buildKey("test", "test2")
	c.Equal("test:test2", key)

	key = buildKey("test", "test2", "test3")
	c.Equal("test:test2:test3", key)
}
