package main

import (
	"fmt"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLogoutHandlerSuccess(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	usersMock := users.NewDynamoMock()
	redisMock := cache.NewRedisCacheMock()

	request := &requestLogoutHandler{
		log:                 logMock,
		userRepo:            users.NewRepository(usersMock),
		invalidTokenManager: redisMock,
	}

	apigwRequest := &apigateway.Request{
		Headers: map[string]string{
			"Cookie": fmt.Sprintf("%s=%s", refreshTokenCookieName, users.DummyToken),
		},
	}

	response, err := request.processLogout(apigwRequest)
	c.NoError(err)
	c.NotEmpty(response.Headers["Set-Cookie"])
	c.Contains(response.Headers["Set-Cookie"], fmt.Sprintf(`%s=;`, refreshTokenCookieName))
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestLogoutHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	usersMock := users.NewDynamoMock()
	redisMock := cache.NewRedisCacheMock()

	request := &requestLogoutHandler{
		log:                 logMock,
		userRepo:            users.NewRepository(usersMock),
		invalidTokenManager: redisMock,
	}

	apigwRequest := &apigateway.Request{
		Headers: map[string]string{
			"Cookie": fmt.Sprintf("%s=%s", refreshTokenCookieName, users.DummyToken),
		},
	}

	t.Run("Cookie header not found", func(t *testing.T) {
		apigwRequest.Headers = map[string]string{}

		response, err := request.processLogout(apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "getting_refresh_token_cookie_failed")
		c.Empty(response.Headers["Set-Cookie"])
		logMock.Output.Reset()
	})

	t.Run("Invalid token length", func(t *testing.T) {
		apigwRequest.Headers["Cookie"] = fmt.Sprintf("%s=%s", refreshTokenCookieName, "a")

		response, err := request.processLogout(apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		logMock.Output.Reset()
	})
}
