package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"net/http"
	"testing"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/stretchr/testify/require"
)

func TestSignUpHandlerFailed(t *testing.T) {
	c := require.New(t)

	body := signUpBody{
		FullName:    "Joel",
		Credentials: &Credentials{"test@gmail.com", "1234"},
	}

	ctx := context.Background()

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	usersMock := users.NewDynamoMock()

	request := &requestSignUpHandler{
		userRepo:         usersMock,
		idempotenceCache: cache.NewRedisCacheMock(),
		secretsManager:   secrets.NewSecretMock(),
	}

	t.Run("Existing user error", func(t *testing.T) {
		usersMock.ActivateForceFailure(models.ErrExistingUser)
		defer usersMock.DeactivateForceFailure()

		apigwRequest := &apigateway.Request{
			Body: jsonBody,
			Headers: map[string]string{
				"Idempotency-Key": "123",
			},
		}

		response, err := request.processSignUp(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, "This account already exists")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest := &apigateway.Request{Body: "}"}

		response, err := request.processSignUp(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	type testCase struct {
		description string
		expectedErr string
		body        signUpBody
	}

	testCases := []testCase{
		{
			"Missing email error",
			"Missing username",
			signUpBody{"", &Credentials{"", "1234"}},
		},
		{
			"Invalid email error",
			"Invalid email",
			signUpBody{"1234", &Credentials{"1234", "1234"}},
		},
		{
			"Missing password error",
			"Missing password",
			signUpBody{"test@gmail.com", &Credentials{"test@gmail.com", ""}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c := require.New(t)

			jsonBody, err = bodyToJSONString(tc.body)
			c.Nil(err)

			apigwRequest := &apigateway.Request{
				Body: jsonBody,
				Headers: map[string]string{
					"Idempotency-Key": "123",
				},
			}

			response, err := request.processSignUp(ctx, apigwRequest)
			c.Nil(err)
			c.Equal(http.StatusBadRequest, response.StatusCode)
			c.Contains(response.Body, tc.expectedErr)
		})
	}
}
