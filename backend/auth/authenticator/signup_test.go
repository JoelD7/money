package main

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"net/http"
	"testing"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/stretchr/testify/require"
)

func TestSignUpHandler(t *testing.T) {
	c := require.New(t)

	body := signUpBody{
		FullName:    "Joel",
		Credentials: &Credentials{"test@gmail.com", "1234"},
	}

	usersMock := users.NewDynamoMock()

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &requestSignUpHandler{
		userRepo: usersMock,
		log:      logger.NewLoggerMock(nil),
	}

	apigwRequest := &apigateway.Request{Body: jsonBody}

	response, err := request.processSignUp(apigwRequest)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func TestSignUpHandlerFailed(t *testing.T) {
	c := require.New(t)

	body := signUpBody{
		FullName:    "Joel",
		Credentials: &Credentials{"test@gmail.com", "1234"},
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	usersMock := users.NewDynamoMock()

	request := &requestSignUpHandler{
		userRepo: usersMock,
		log:      logger.NewLoggerMock(nil),
	}

	t.Run("Existing user error", func(t *testing.T) {
		usersMock.ActivateForceFailure(models.ErrExistingUser)
		defer usersMock.DeactivateForceFailure()

		apigwRequest := &apigateway.Request{Body: jsonBody}

		response, err := request.processSignUp(apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Equal(models.ErrExistingUser.Error(), response.Body)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest := &apigateway.Request{Body: "}"}

		response, err := signUpHandler(apigwRequest)
		c.NotNil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Equal(apigateway.ErrInternalError.Message, response.Body)
	})

	type testCase struct {
		description string
		expectedErr error
		body        signUpBody
	}

	testCases := []testCase{
		{
			"Missing email error",
			models.ErrMissingEmail,
			signUpBody{"", &Credentials{"", "1234"}},
		},
		{
			"Invalid email error",
			models.ErrInvalidEmail,
			signUpBody{"1234", &Credentials{"1234", "1234"}},
		},
		{
			"Missing password error",
			models.ErrMissingPassword,
			signUpBody{"test@gmail.com", &Credentials{"test@gmail.com", ""}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c := require.New(t)

			jsonBody, err = bodyToJSONString(tc.body)
			c.Nil(err)

			apigwRequest := &apigateway.Request{Body: jsonBody}

			response, err := request.processSignUp(apigwRequest)
			c.Nil(err)
			c.Equal(http.StatusBadRequest, response.StatusCode)
			c.Equal(tc.expectedErr.Error(), response.Body)
		})
	}
}
