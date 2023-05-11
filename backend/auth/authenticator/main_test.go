package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/shared/logger"
	restMock "github.com/JoelD7/money/backend/shared/restclient/mocks"
	secretsMock "github.com/JoelD7/money/backend/shared/secrets/mocks"
	storageInvalidToken "github.com/JoelD7/money/backend/storage/invalid_token"
	storagePerson "github.com/JoelD7/money/backend/storage/person"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

var secretMock *secretsMock.MockSecret

func init() {
	logger.InitLoggerMock()

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

	err := mockRestClientGetFromFile("samples/jwks_response.json")
	c.Nil(err)

	response, err := jwksHandler(&events.APIGatewayProxyRequest{})
	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal(expectedJWKS, response.Body)
}

func TestRefreshTokenHandler(t *testing.T) {
	c := require.New(t)

	_ = storagePerson.InitDynamoMock()

	person := storagePerson.GetMockedPerson()

	request, err := dummyAPIGatewayProxyRequest()
	c.Nil(err)

	request.Headers["Cookie"] = refreshTokenCookieName + "=" + person.RefreshToken

	response, err := refreshTokenHandler(&request)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
	c.NotEmpty(response.Body)
}

func TestRefreshTokenHandlerFailed(t *testing.T) {
	c := require.New(t)

	dummyRequest, err := dummyAPIGatewayProxyRequest()
	c.Nil(err)

	t.Run("Invalid request body", func(t *testing.T) {
		request := &events.APIGatewayProxyRequest{Body: "}"}

		response, err := refreshTokenHandler(request)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Refresh token leaked", func(t *testing.T) {
		request := dummyRequest

		request.Headers = map[string]string{}
		request.Headers["Cookie"] = refreshTokenCookieName + "=previous token"

		response, err := refreshTokenHandler(&request)
		c.Nil(err)
		c.Equal(http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Person not found", func(t *testing.T) {
		personMock := storagePerson.InitDynamoMock()

		personMock.ActivateForceFailure(storagePerson.ErrNotFound)
		defer personMock.DeactivateForceFailure()

		request := dummyRequest

		response, err := refreshTokenHandler(&request)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Invalid person access token", func(t *testing.T) {
		personMock := storagePerson.InitDynamoMock()

		person := storagePerson.GetMockedPerson()

		person.AccessToken = "invalid token"

		err = personMock.MockQueryFromSource(person)
		c.Nil(err)

		request := dummyRequest

		request.Headers["Cookie"] = refreshTokenCookieName + "=" + person.PreviousRefreshToken

		response, err := refreshTokenHandler(&request)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Invalid person refresh token", func(t *testing.T) {
		personMock := storagePerson.InitDynamoMock()

		person := storagePerson.GetMockedPerson()

		person.RefreshToken = "invalid token"

		err = personMock.MockQueryFromSource(person)
		c.Nil(err)

		request := dummyRequest

		request.Headers["Cookie"] = refreshTokenCookieName + "=" + "random"

		response, err := refreshTokenHandler(&request)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Token invalidation failed", func(t *testing.T) {
		_ = storagePerson.InitDynamoMock()
		person := storagePerson.GetMockedPerson()

		itMock := storageInvalidToken.InitDynamoMock()

		request := dummyRequest

		request.Headers["Cookie"] = refreshTokenCookieName + "=" + person.PreviousRefreshToken

		errCustomError := errors.New("custom error")

		itMock.ActivateForceFailure(errCustomError)
		defer itMock.DeactivateForceFailure()

		response, err := refreshTokenHandler(&request)
		c.EqualError(errCustomError, err.Error())
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})
}

func bodyToJSONString(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func mockRestClientGetFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	r := io.NopCloser(bytes.NewReader(data))

	restMock.GetFunction = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	return nil
}

func dummyAPIGatewayProxyRequest() (events.APIGatewayProxyRequest, error) {
	body := Credentials{
		Email: "test@gmail.com",
	}

	jsonBody, err := bodyToJSONString(body)
	if err != nil {
		return events.APIGatewayProxyRequest{}, err
	}

	return events.APIGatewayProxyRequest{
		Body:    jsonBody,
		Headers: map[string]string{},
	}, nil
}

//func getDummyPerson() (*models.Person, error) {
//	dummyToken, err := getDummyToken()
//	if err != nil {
//		return nil, err
//	}
//
//	return &models.Person{
//		FullName:             "Joel",
//		Email:                "test@gmail.com",
//		Password:             "$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe",
//		PreviousRefreshToken: "previous token",
//		AccessToken:          dummyToken,
//		RefreshToken:         dummyToken,
//	}, nil
//}
//
//func getDummyToken() (string, error) {
//	pld := &models.JWTPayload{
//		Payload: &jwt.Payload{
//			Subject:        "John Doe",
//			ExpirationTime: jwt.NumericDate(time.Now().Add(time.Hour * 1)),
//		},
//	}
//
//	payload, err := json.Marshal(pld)
//	if err != nil {
//		return "", err
//	}
//
//	encodedPayload := make([]byte, base64.RawURLEncoding.EncodedLen(len(payload)))
//	base64.RawURLEncoding.Encode(encodedPayload, payload)
//
//	return "random." + string(encodedPayload) + ".random", nil
//}
