package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/stretchr/testify/require"
)

func TestTokenHandler(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	secretMock := secrets.NewSecretMock()
	ctx := context.Background()

	secretMock.RegisterResponder(privateSecretName, func(ctx context.Context, name string) (string, error) {
		return "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA5l+M6MGnS6K8SNXUIqOGaaH/IO7NcBxwQJVd4X6uUcLHfdhy\nNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4PoIUikiFnuCmkVDxcnl65/jv4DQt\nDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFyntaOREIMTaLQ3F72aDMk0ytVFu0c\nZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8/ZKklkKTDnOmGlbVOKTvujH6fTJu\nQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqOWTuEpsvMNRubWoLt22c9f4PXaGDw\nqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98HQIDAQABAoIBAQC8OffcuVVihC2I\n6UUxpCCPsG/PTa6HWoURD7msI6B0Z0wt86qkPCH0xlxufhxt/wk4GvIiEqm2P5YG\n7l0JUGh3EjuMHOMoQ+90rgkI+l7EqG3950OjtQvHP4aoF0BDlgZAv4h3FUl5dJsw\neow2mZfoVe/Zr4lz6YLze5ei3Z7J9YGjj62j7QGbKgbwPLqnqnNrUQqM4T0V9SaJ\nCE6sDxYo8M8kE2yqgiIvsA1D92u4AMchcdnjREBy/ogCRXzvuZQADC4st8UE4aFD\nmDNKwIbprymSa7atjSMz+lfWWBnuuzFmsf+72gXJVRRpmbm7onBxDHcJlqk0fMjv\nm+zuow4hAoGBAPyGYveulOLbwA+88DHBJ5GVKZNJHRGFsMWpysCC45OC5eKs9yrP\nnI6/0mvL1JFXvSbkkDql6qvlbKVcoH80h3ipFwuB4j8KbLXs/LiFMF0yg4IpAoq5\nGIp1RK6VBsv4OvP8vmkJxzakgRn5C7WyswNqJHGW108l7pmYEbOpfohZAoGBAOmL\nIE3ZbttXLOLEcHICcySpNdDTIWaiNqMtyjLLK29Ic8KDpXI1hfYzAw+7Q6y2QyEG\nl9l5IpkY6Wt29LMuWUMH6fnS1H/JeQOSnkT2y8PXSN95QKb2HtbP+ujqSdG12EPs\nILag0918ezcttGLszqWfipSZuSo2ZQ6b0A+uaxllAoGAECkBeFxBxurNNbSfom97\n+sMS8AwDwjVOBLhC82Ls8Wm1EHaFMsYqfLAl5SQcLFjzD+QcnsQzamC6PTLaSomw\nCba4dNIRCnu+TT4nRh+v4qby54d8VChYO7QZexqqXq86BpcsEEjB6OtKH8FiUHRp\nJFTMlEBU8wm4ZTfoGhlEsbECgYAkNJ1ddEfrWShsP2fvRNH07QaaySB0eNFfmsmt\n9jFVnzXTAfW0LvgFowLmfXGQZPEjPZJs9IqYkXQeZOKqpJTR/3gWcsjexq0sEJ7Y\nsioEwmtZucJ8H8vIIZYUZb3r9PUCEqk/ps8xlwrDEyLT80JWCtXBE9PQ533jNeSb\nib6wwQKBgQDkCKsHfxv/z+YgdMe3mUCSZi2gNttPQczjeUSYAxYITj/OJ1TfMuk4\n8gVdOcusHynFH3jEpnA8fqdpZpmhH/sAKPuQl/vwBCefVyBO5LkM14gxEIf9eq69\n7QVBd9ep1cN/5yYcJUJAcpjBxcbR8rXYowLtsYaGsC7G5tMlW8rJTg==\n-----END RSA PRIVATE KEY-----", nil
	})

	request := &requestTokenHandler{
		log:                 logger.NewLoggerMock(nil),
		secretsManager:      secretMock,
		userRepo:            usersMock,
		invalidTokenManager: cache.NewRedisCacheMock(),
	}

	apigwRequest, err := dummyAPIGatewayProxyRequest()
	c.Nil(err)

	apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken

	response, err := request.processToken(ctx, apigwRequest)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
	c.Contains(response.Body, "access_token")
	c.NotNil(response.Headers["Set-Cookie"])
	c.Contains(response.Headers["Set-Cookie"], refreshTokenCookieName)
}

func TestTokenHandlerFailed(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()

	dummyApigwRequest, err := dummyAPIGatewayProxyRequest()
	dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken
	c.Nil(err)

	usersMock := users.NewDynamoMock()
	secretMock := secrets.NewSecretMock()
	redisMock := cache.NewRedisCacheMock()
	logMock := logger.NewLoggerMock(nil)

	request := &requestTokenHandler{
		log:                 logMock,
		secretsManager:      secretMock,
		userRepo:            usersMock,
		invalidTokenManager: redisMock,
	}

	t.Run("Invalid token", func(t *testing.T) {
		dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "="

		response, err := request.processToken(ctx, dummyApigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "getting_refresh_token_cookie_failed")

		dummyApigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=header.payload.signature"
		response, err = request.processToken(ctx, dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_refresh_token_payload_failed")
	})

	t.Run("Refresh token leaked", func(t *testing.T) {
		apigwRequest := dummyApigwRequest

		apigwRequest.Headers = map[string]string{}
		apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyPreviousToken

		response, err := request.processToken(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "refresh_token_validation_failed")
	})

	t.Run("Token invalidation failed", func(t *testing.T) {
		dummyErr := errors.New("dummy error")

		redisMock.ActivateForceFailure(dummyErr)
		defer redisMock.DeactivateForceFailure()

		apigwRequest := dummyApigwRequest

		apigwRequest.Headers = map[string]string{}
		apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyPreviousToken

		response, err := request.processToken(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "refresh_token_validation_failed")
	})

	t.Run("User not found", func(t *testing.T) {
		usersMock.ActivateForceFailure(models.ErrUserNotFound)
		defer usersMock.DeactivateForceFailure()

		apigwRequest := dummyApigwRequest

		response, err := request.processToken(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "get_user_failed")
	})

	t.Run("Refresh token in cookie not found", func(t *testing.T) {
		dummyApigwRequest.Headers["Cookie"] = ""

		response, err := request.processToken(ctx, dummyApigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "getting_refresh_token_cookie_failed")
	})

	t.Run("Set tokens failed", func(t *testing.T) {
		secretMock.ActivateForceFailure(secrets.SecretsError)
		defer secretMock.DeactivateForceFailure()

		apigwRequest := dummyApigwRequest
		apigwRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + users.DummyToken

		response, err := request.processToken(ctx, dummyApigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Empty(response.Headers["Set-Cookie"])
		c.Contains(logMock.Output.String(), "generate_access_token_failed")
	})
}

func dummyAPIGatewayProxyRequest() (*apigateway.Request, error) {
	body := Credentials{
		Username: "test@gmail.com",
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

func bodyToJSONString(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
