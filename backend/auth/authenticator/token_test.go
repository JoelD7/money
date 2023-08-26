package main

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestTokenHandler(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	secretMock := secrets.NewSecretMock()
	redisRepository := cache.NewRepository(cache.NewRedisCacheMock())

	request := &requestTokenHandler{
		log:            logger.NewLoggerWithHandler("token"),
		secretsManager: secretMock,
		userRepo:       users.NewRepository(usersMock),
		cacheRepo:      redisRepository,
	}

	apigwRequest, err := dummyAPIGatewayProxyRequest()
	c.Nil(err)

	apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken

	response, err := request.processToken(apigwRequest)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
	c.NotEmpty(response.Body)
}

func TestTokenHandlerFailed(t *testing.T) {
	c := require.New(t)

	dummyApigwRequest, err := dummyAPIGatewayProxyRequest()
	dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken
	c.Nil(err)

	usersMock := users.NewDynamoMock()
	secretMock := secrets.NewSecretMock()
	redisRepository := cache.NewRepository(cache.NewRedisCacheMock())

	request := &requestTokenHandler{
		log:            logger.NewLoggerWithHandler("token"),
		secretsManager: secretMock,
		userRepo:       users.NewRepository(usersMock),
		cacheRepo:      redisRepository,
	}

	t.Run("Invalid token", func(t *testing.T) {
		dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "="

		response, err := request.processToken(dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_payload_parse_failed")

		dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=header.payload.signature"
		response, err = request.processToken(dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_payload_parse_failed")
	})

	t.Run("Refresh token leaked", func(t *testing.T) {
		apigwRequest := dummyApigwRequest

		apigwRequest.Headers = map[string]string{}
		apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyPreviousToken

		response, err := request.processToken(apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
		c.Contains(logMock.Output.String(), "refresh_token_validation_failed")
	})

	t.Run("Person not found", func(t *testing.T) {
		usersMock.ActivateForceFailure(models.ErrUserNotFound)
		usersMock.DeactivateForceFailure()

		apigwRequest := dummyApigwRequest

		response, err := request.processToken(apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "fetching_user_from_storage_failed")
	})

	t.Run("Refresh token in cookie not found", func(t *testing.T) {
		dummyApigwRequest.Headers["Cookie"] = ""

		response, err := request.processToken(dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "getting_refresh_token_cookie_failed")

	})

	t.Run("Set tokens failed", func(t *testing.T) {
		secretMock.ActivateForceFailure(secrets.SecretsError)
		defer secretMock.DeactivateForceFailure()

		apigwRequest := dummyApigwRequest
		apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken

		response, err := request.processToken(dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_setting_failed")
	})
}

func dummyAPIGatewayProxyRequest() (*apigateway.Request, error) {
	body := Credentials{
		Email: "test@gmail.com",
	}

	jsonBody, err := bodyToJSONString(body)
	if err != nil {
		return &apigateway.Request{}, err
	}

	return &apigateway.Request{
		Body:    jsonBody,
		Headers: map[string]string{},
	}, nil
}
