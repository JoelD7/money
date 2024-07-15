package env

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestJ(t *testing.T) {
	c := require.New(t)

	data, err := os.ReadFile("../../.env")
	c.Nil(err)
	fmt.Println(string(data))
}
