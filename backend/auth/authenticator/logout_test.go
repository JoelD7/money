package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLogoutHandlerSuccess(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	redisMock := cache.NewRedisCacheMock()
	ctx := context.Background()

	request := &requestLogoutHandler{
		userRepo:            usersMock,
		invalidTokenManager: redisMock,
	}

	apigwRequest := &apigateway.Request{
		Body: `{"username":"test@gmail.com"}`,
	}

	response, err := request.processLogout(ctx, apigwRequest)
	c.NoError(err)
	c.NotEmpty(response.Headers["Set-Cookie"])
	c.Contains(response.Headers["Set-Cookie"], fmt.Sprintf(`%s=;`, refreshTokenCookieName))
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestLogoutHandlerFailed(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	usersMock := users.NewDynamoMock()
	redisMock := cache.NewRedisCacheMock()

	request := &requestLogoutHandler{
		userRepo:            usersMock,
		invalidTokenManager: redisMock,
	}

	t.Run("Empty username", func(t *testing.T) {
		apigwRequest := &apigateway.Request{
			Body: `{"username":""}`,
		}

		response, err := request.processLogout(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
	})
}
