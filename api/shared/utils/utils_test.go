package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateDynamoID(t *testing.T) {
	c := require.New(t)

	set := make(map[string]struct{})
	maxIterations := 10000
	var id string

	for i := 0; i < maxIterations; i++ {
		id = GenerateDynamoID("rd")

		_, ok := set[id]
		c.False(ok)

		set[id] = struct{}{}
	}
}
