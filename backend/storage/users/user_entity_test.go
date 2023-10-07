package users

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToUserEntity(t *testing.T) {
	c := require.New(t)

	user := GetDummyUser()

	userEntity := toUserEntity(user)
	c.Len(userEntity.Categories, 3)
}
