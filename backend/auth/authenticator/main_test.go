package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	secretsMock "github.com/JoelD7/money/backend/shared/secrets/mocks"
	"github.com/JoelD7/money/backend/storage/invalidtoken"
	storagePerson "github.com/JoelD7/money/backend/storage/person"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var (
	secretMock *secretsMock.MockSecret
	logMock    *logger.LogMock

	logBuffer bytes.Buffer
)

func init() {
	logMock = logger.InitLoggerMock(logBuffer)

	secretMock = secretsMock.InitSecretMock()

	secretMock.RegisterResponder(publicSecretName, func(ctx context.Context, name string) (string, error) {
		return "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5l+M6MGnS6K8SNXUIqOG\naaH/IO7NcBxwQJVd4X6uUcLHfdhyNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4\nPoIUikiFnuCmkVDxcnl65/jv4DQtDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFy\nntaOREIMTaLQ3F72aDMk0ytVFu0cZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8\n/ZKklkKTDnOmGlbVOKTvujH6fTJuQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqO\nWTuEpsvMNRubWoLt22c9f4PXaGDwqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98\nHQIDAQAB\n-----END PUBLIC KEY-----", nil
	})

	secretMock.RegisterResponder(privateSecretName, func(ctx context.Context, name string) (string, error) {
		return "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA5l+M6MGnS6K8SNXUIqOGaaH/IO7NcBxwQJVd4X6uUcLHfdhy\nNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4PoIUikiFnuCmkVDxcnl65/jv4DQt\nDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFyntaOREIMTaLQ3F72aDMk0ytVFu0c\nZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8/ZKklkKTDnOmGlbVOKTvujH6fTJu\nQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqOWTuEpsvMNRubWoLt22c9f4PXaGDw\nqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98HQIDAQABAoIBAQC8OffcuVVihC2I\n6UUxpCCPsG/PTa6HWoURD7msI6B0Z0wt86qkPCH0xlxufhxt/wk4GvIiEqm2P5YG\n7l0JUGh3EjuMHOMoQ+90rgkI+l7EqG3950OjtQvHP4aoF0BDlgZAv4h3FUl5dJsw\neow2mZfoVe/Zr4lz6YLze5ei3Z7J9YGjj62j7QGbKgbwPLqnqnNrUQqM4T0V9SaJ\nCE6sDxYo8M8kE2yqgiIvsA1D92u4AMchcdnjREBy/ogCRXzvuZQADC4st8UE4aFD\nmDNKwIbprymSa7atjSMz+lfWWBnuuzFmsf+72gXJVRRpmbm7onBxDHcJlqk0fMjv\nm+zuow4hAoGBAPyGYveulOLbwA+88DHBJ5GVKZNJHRGFsMWpysCC45OC5eKs9yrP\nnI6/0mvL1JFXvSbkkDql6qvlbKVcoH80h3ipFwuB4j8KbLXs/LiFMF0yg4IpAoq5\nGIp1RK6VBsv4OvP8vmkJxzakgRn5C7WyswNqJHGW108l7pmYEbOpfohZAoGBAOmL\nIE3ZbttXLOLEcHICcySpNdDTIWaiNqMtyjLLK29Ic8KDpXI1hfYzAw+7Q6y2QyEG\nl9l5IpkY6Wt29LMuWUMH6fnS1H/JeQOSnkT2y8PXSN95QKb2HtbP+ujqSdG12EPs\nILag0918ezcttGLszqWfipSZuSo2ZQ6b0A+uaxllAoGAECkBeFxBxurNNbSfom97\n+sMS8AwDwjVOBLhC82Ls8Wm1EHaFMsYqfLAl5SQcLFjzD+QcnsQzamC6PTLaSomw\nCba4dNIRCnu+TT4nRh+v4qby54d8VChYO7QZexqqXq86BpcsEEjB6OtKH8FiUHRp\nJFTMlEBU8wm4ZTfoGhlEsbECgYAkNJ1ddEfrWShsP2fvRNH07QaaySB0eNFfmsmt\n9jFVnzXTAfW0LvgFowLmfXGQZPEjPZJs9IqYkXQeZOKqpJTR/3gWcsjexq0sEJ7Y\nsioEwmtZucJ8H8vIIZYUZb3r9PUCEqk/ps8xlwrDEyLT80JWCtXBE9PQ533jNeSb\nib6wwQKBgQDkCKsHfxv/z+YgdMe3mUCSZi2gNttPQczjeUSYAxYITj/OJ1TfMuk4\n8gVdOcusHynFH3jEpnA8fqdpZpmhH/sAKPuQl/vwBCefVyBO5LkM14gxEIf9eq69\n7QVBd9ep1cN/5yYcJUJAcpjBxcbR8rXYowLtsYaGsC7G5tMlW8rJTg==\n-----END RSA PRIVATE KEY-----", nil
	})

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
		return "123", nil
	})
}

func TestLoginHandler(t *testing.T) {
	c := require.New(t)

	body := Credentials{
		Email:    "test@gmail.com",
		Password: "1234",
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	response, err := logInHandler(request)
	c.Equal(http.StatusOK, response.StatusCode)
	c.NotNil(response.Headers["Set-Cookie"])
	c.Contains(response.Headers["Set-Cookie"], "refresh_token")
	c.Contains(response.Body, "access_token")
}

func TestLoginHandlerFailed(t *testing.T) {
	c := require.New(t)

	body := Credentials{
		Email:    "test@gmail.com",
		Password: "1234",
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	secretMock.ActivateForceFailure(secretsMock.SecretsError)
	defer secretMock.DeactivateForceFailure()

	response, err := logInHandler(request)
	c.Nil(err)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal(http.StatusText(http.StatusInternalServerError), response.Body)

	personMock := storagePerson.InitDynamoMock()
	personMock.ActivateForceFailure(storagePerson.ErrNotFound)

	response, err = logInHandler(request)
	c.Nil(err)
	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal(storagePerson.ErrNotFound.Error(), response.Body)
	personMock.DeactivateForceFailure()

	request.Body = "a"
	response, err = logInHandler(request)
	c.Nil(err)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal(http.StatusText(http.StatusInternalServerError), response.Body)

	type testCase struct {
		description string
		expectedErr string
		body        Credentials
	}

	testCases := []testCase{
		{
			"Wrong credentials",
			errWrongCredentials.Error(),
			Credentials{"test@gmail.com", "random"},
		},
		{
			"Missing email error",
			errMissingEmail.Error(),
			Credentials{"", "1234"},
		},
		{
			"Invalid email error",
			errInvalidEmail.Error(),
			Credentials{"1234", "1234"},
		},
		{
			"Missing password error",
			errMissingPassword.Error(),
			Credentials{"test@gmail.com", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c := require.New(t)

			jsonBody, err = bodyToJSONString(tc.body)
			c.Nil(err)

			request.Body = jsonBody

			response, err = logInHandler(request)
			c.Nil(err)
			c.Equal(http.StatusBadRequest, response.StatusCode)
			c.Equal(tc.expectedErr, response.Body)
		})
	}
}

func TestSignUpHandler(t *testing.T) {
	c := require.New(t)

	body := signUpBody{
		FullName:    "Joel",
		Credentials: &Credentials{"test@gmail.com", "1234"},
	}

	personMock := storagePerson.InitDynamoMock()

	personMock.EmptyTable()

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	response, err := signUpHandler(request)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestSignUpHandlerFailed(t *testing.T) {
	c := require.New(t)

	body := signUpBody{
		FullName:    "Joel",
		Credentials: &Credentials{"test@gmail.com", "1234"},
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	personMock := storagePerson.InitDynamoMock()

	personMock.ActivateForceFailure(storagePerson.ErrExistingUser)
	defer personMock.DeactivateForceFailure()

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	response, err := signUpHandler(request)
	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal(storagePerson.ErrExistingUser.Error(), response.Body)

	request = &events.APIGatewayProxyRequest{Body: "}"}

	response, err = signUpHandler(request)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal(http.StatusText(http.StatusInternalServerError), response.Body)

	type testCase struct {
		description string
		expectedErr string
		body        signUpBody
	}

	testCases := []testCase{
		{
			"Missing email error",
			errMissingEmail.Error(),
			signUpBody{"", &Credentials{"", "1234"}},
		},
		{
			"Invalid email error",
			errInvalidEmail.Error(),
			signUpBody{"1234", &Credentials{"1234", "1234"}},
		},
		{
			"Missing password error",
			errMissingPassword.Error(),
			signUpBody{"test@gmail.com", &Credentials{"test@gmail.com", ""}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c := require.New(t)

			jsonBody, err = bodyToJSONString(tc.body)
			c.Nil(err)

			request.Body = jsonBody

			response, err = signUpHandler(request)
			c.Nil(err)
			c.Equal(http.StatusBadRequest, response.StatusCode)
			c.Equal(tc.expectedErr, response.Body)
		})
	}
}

func TestJWTHandler(t *testing.T) {
	c := require.New(t)

	expectedJWKS := `{"keys":[{"kty":"RSA","kid":"123","use":"sig","n":"5l-M6MGnS6K8SNXUIqOGaaH_IO7NcBxwQJVd4X6uUcLHfdhyNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8_k4PoIUikiFnuCmkVDxcnl65_jv4DQtDL6GGqoLcYo2ENldfj8uDo09CmYS_DKuJxFyntaOREIMTaLQ3F72aDMk0ytVFu0cZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK_D8_ZKklkKTDnOmGlbVOKTvujH6fTJuQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqOWTuEpsvMNRubWoLt22c9f4PXaGDwqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98HQ","e":"AQAB"}]}`

	err := restclient.AddMockedResponseFromFile("samples/jwks_response.json", accessTokenIssuer+"/auth/jwks", restclient.MethodGET)
	c.Nil(err)

	response, err := jwksHandler(&events.APIGatewayProxyRequest{})
	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal(expectedJWKS, response.Body)
}

func TestTokenHandler(t *testing.T) {
	c := require.New(t)

	_ = storagePerson.InitDynamoMock()

	request, err := dummyAPIGatewayProxyRequest()
	c.Nil(err)

	request.Headers["Cookie"] = refreshTokenCookieName + "=" + storagePerson.DummyToken

	response, err := tokenHandler(request)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
	c.NotEmpty(response.Body)
}

func TestTokenHandlerFailed(t *testing.T) {
	c := require.New(t)

	dummyRequest, err := dummyAPIGatewayProxyRequest()
	dummyRequest.Headers["Cookie"] = refreshTokenCookieName + "=" + storagePerson.DummyToken
	c.Nil(err)

	t.Run("Invalid token", func(t *testing.T) {
		dummyRequest.Headers["Cookie"] = refreshTokenCookieName + "="

		response, err := tokenHandler(dummyRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_payload_parse_failed")

		dummyRequest.Headers["Cookie"] = refreshTokenCookieName + "=header.payload.signature"
		response, err = tokenHandler(dummyRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_payload_parse_failed")
	})

	t.Run("Refresh token leaked", func(t *testing.T) {
		request := dummyRequest

		request.Headers = map[string]string{}
		request.Headers["Cookie"] = refreshTokenCookieName + "=" + storagePerson.DummyPreviousToken

		response, err := tokenHandler(request)
		c.Nil(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
		c.Contains(logMock.Output.String(), "invalid_refresh_token")
	})

	t.Run("Person not found", func(t *testing.T) {
		personMock := storagePerson.InitDynamoMock()

		personMock.ActivateForceFailure(storagePerson.ErrNotFound)
		defer personMock.DeactivateForceFailure()

		request := dummyRequest

		response, err := tokenHandler(request)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "fetching_user_from_storage_failed")
	})

	t.Run("Refresh token in cookie not found", func(t *testing.T) {
		_ = storagePerson.InitDynamoMock()

		dummyRequest.Headers["Cookie"] = ""

		response, err := tokenHandler(dummyRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "getting_refresh_token_cookie_failed")

	})

	t.Run("Token invalidation failed", func(t *testing.T) {
		_ = storagePerson.InitDynamoMock()
		person := storagePerson.GetMockedPerson()

		itMock := invalidtoken.InitDynamoMock()

		request := dummyRequest

		request.Headers["Cookie"] = refreshTokenCookieName + "=" + person.PreviousRefreshToken

		errCustomError := errors.New("custom error")

		itMock.ActivateForceFailure(errCustomError)
		defer itMock.DeactivateForceFailure()

		response, err := tokenHandler(request)
		c.EqualError(errCustomError, err.Error())
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "access_token_invalidation_failed")
	})

	t.Run("Set tokens failed", func(t *testing.T) {
		_ = storagePerson.InitDynamoMock()
		_ = invalidtoken.InitDynamoMock()

		sMock := secretsMock.InitSecretMock()

		sMock.ActivateForceFailure(secretsMock.SecretsError)
		defer sMock.DeactivateForceFailure()

		request := dummyRequest
		request.Headers["Cookie"] = refreshTokenCookieName + "=" + storagePerson.DummyToken

		response, err := tokenHandler(dummyRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "token_setting_failed")
	})
}

func bodyToJSONString(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func dummyAPIGatewayProxyRequest() (*events.APIGatewayProxyRequest, error) {
	body := Credentials{
		Email: "test@gmail.com",
	}

	jsonBody, err := bodyToJSONString(body)
	if err != nil {
		return &events.APIGatewayProxyRequest{}, err
	}

	return &events.APIGatewayProxyRequest{
		Body:    jsonBody,
		Headers: map[string]string{},
	}, nil
}
