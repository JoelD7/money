package main

import (
	"bytes"
	"context"
	"encoding/json"
	restMock "github.com/JoelD7/money/api/shared/restclient/mocks"
	secretsMock "github.com/JoelD7/money/api/shared/secrets/mocks"
	"github.com/JoelD7/money/api/storage"
	"github.com/JoelD7/money/api/storage/mocks"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func init() {
	mocks.InitDynamoMock()
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

}

func TestLoginHandlerFailed(t *testing.T) {
	c := require.New(t)

	body := Credentials{
		Email:    "test@gmail.com",
		Password: "random",
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	mocks.ForceNotFound = true

	response, err := logInHandler(request)
	c.Nil(err)
	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal(mocks.ErrForceNotFound.Error(), response.Body)

	mocks.ForceNotFound = false

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

	mocks.ForceUserExists = true

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	response, err := signUpHandler(request)
	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal(storage.ErrExistingUser.Error(), response.Body)
	mocks.ForceUserExists = false

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

			response, err = logInHandler(request)
			c.Nil(err)
			c.Equal(http.StatusBadRequest, response.StatusCode)
			c.Equal(tc.expectedErr, response.Body)
		})
	}
}

func TestJWTHandler(t *testing.T) {
	c := require.New(t)

	expectedJWKS := `{"keys":[{"kty":"RSA","kid":"123","use":"sig","n":"qGtV1QpRQ6he8z3l64alazzW4dBnfOUF_J1EDTP7i8DJPhlFFE1Mn-zTZN_-jGgMjhHUHG3AUfv2khUR0Bi4T0DnQlSrlW_TcT2747AEu8qTAgXagUDy3YhwGiqsBy-S_fv0zGgVbRLeqNKnYqEAgQDhX7EbIyx9ke00jM6tbEeguOtCp6VoslRN3rM_yqi0xKHOxIoTbTedmg-cBqqmMZYyanLnAuzjYrrieW-23O_YkV0tbTJjhL_XJXeBze0C8Iltcvfaxhlxd_jpm28gO01n91PKwg-YhwPhYIpxlzrKps0mo6iAhNvsDNGFha_8UiZ-bJa5F7xk3LArTrbvbQ","e":"AQAB"}]}`

	err := mockRestClientGetFromFile("samples/jwks_response.json")
	c.Nil(err)

	secretMock := secretsMock.InitSecretMock()

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("123"),
		}, nil
	})

	secretMock.RegisterResponder(publicSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqGtV1QpRQ6he8z3l64al\nazzW4dBnfOUF/J1EDTP7i8DJPhlFFE1Mn+zTZN/+jGgMjhHUHG3AUfv2khUR0Bi4\nT0DnQlSrlW/TcT2747AEu8qTAgXagUDy3YhwGiqsBy+S/fv0zGgVbRLeqNKnYqEA\ngQDhX7EbIyx9ke00jM6tbEeguOtCp6VoslRN3rM/yqi0xKHOxIoTbTedmg+cBqqm\nMZYyanLnAuzjYrrieW+23O/YkV0tbTJjhL/XJXeBze0C8Iltcvfaxhlxd/jpm28g\nO01n91PKwg+YhwPhYIpxlzrKps0mo6iAhNvsDNGFha/8UiZ+bJa5F7xk3LArTrbv\nbQIDAQAB\n-----END PUBLIC KEY-----"),
		}, nil
	})

	response, err := jwksHandler(&events.APIGatewayProxyRequest{})
	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal(expectedJWKS, response.Body)
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
