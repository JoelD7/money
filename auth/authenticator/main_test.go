package main

import (
	"encoding/json"
	storage "github.com/JoelD7/money/api/storage/person"
	"github.com/JoelD7/money/api/storage/person/mocks"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
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

func bodyToJSONString(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
