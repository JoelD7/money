package main

import (
	"bytes"
	"context"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	secretsMock "github.com/JoelD7/money/backend/shared/secrets/mocks"
	"github.com/JoelD7/money/backend/storage/invalidtoken"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	authToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6InJlYWQgd3JpdGUiLCJpc3MiOiJodHRwczovLzM4cXNscGU4ZDkuZXhlY3V0ZS1hcGkudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vc3RhZ2luZyIsInN1YiI6InRlc3RAZ21haWwuY29tIiwiYXVkIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6MzAwMCIsImV4cCI6MTcwODI5OTA4OCwibmJmIjoxNjc3MTk2ODg4LCJpYXQiOjE2NzcxOTUwODh9.S_wnwVHTs_-T9zOkIFVIblfYYZ338kgUDclRi5nzgzLxzfqo_jrxKYXwLVeVkRNq1etO4B2RmyFPsLVHpC4cGS_Kr093eOzdWta0F8nj_hbTK2ZtuNP88X8oKaDadyCbXFw3M6dxm0la9kf20CZRxFsbtJ0MqPBqW9lp3B_XRz_pTAqMQnVbyfmbQBZiGBKpK5Ur1g043YAP5B2cd2C0ARGyyWw1UzXJBZbM_8KUFLUtndjZIn_uF3z8fLaH4hrnN3Gz_CnRIhgb6kbAWJ2OWsSJb4l15vgzdw2GvOWHU7MHqX6VoIwPVUzFTMDHkzfjDjhnKdWDj2bL-I-XXvZgSg"

	secretMock *secretsMock.MockSecret
	logMock    *logger.LogMock
	logBuffer  bytes.Buffer
)

func init() {
	restclient.InitMockClient()
	logMock = logger.InitLoggerMock(logBuffer)

	secretMock = secretsMock.InitSecretMock()

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
		return "123", nil
	})
}

func TestHandleRequest(t *testing.T) {
	c := require.New(t)

	_ = invalidtoken.InitDynamoMock()

	event := dummyHandlerEvent()

	err := restclient.AddMockedResponseFromFile("samples/jwks_response.json", jwtIssuer+"/auth/jwks", restclient.MethodGET)
	c.Nil(err)

	_, err = handleRequest(context.Background(), event)
	c.Nil(err)
}

func TestHandlerError(t *testing.T) {
	c := require.New(t)

	event := dummyHandlerEvent()
	authToken := event.AuthorizationToken

	t.Run("Invalid token length", func(t *testing.T) {
		event.AuthorizationToken = "dummy"

		_, err := handleRequest(context.Background(), event)
		c.ErrorIs(err, errUnauthorized)
		c.Contains(logMock.Output.String(), "invalid_token_length_detected")
		logMock.Output.Reset()
	})

	t.Run("Payload decoding failed", func(t *testing.T) {
		event.AuthorizationToken = "Bearer dummy.dummy.token"

		_, err := handleRequest(context.Background(), event)
		c.ErrorIs(err, errUnauthorized)
		c.Contains(logMock.Output.String(), "payload_decoding_failed")
		logMock.Output.Reset()
	})

	t.Run("Signing key not found", func(t *testing.T) {
		secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
			return "456", nil
		})

		err := restclient.AddMockedResponseFromFile("samples/jwks_response.json", jwtIssuer+"/auth/jwks", restclient.MethodGET)
		c.Nil(err)

		event.AuthorizationToken = authToken
		response, err := handleRequest(context.Background(), event)
		c.Nil(err)
		c.NotNil(response.Context["stringKey"])
		c.Equal(errSigningKeyNotFound.Error(), response.Context["stringKey"])
		c.Contains(logMock.Output.String(), "signing_key_not_found")
		logMock.Output.Reset()
	})

	t.Run("Getting public key failed", func(t *testing.T) {
		secretMock.ActivateForceFailure(secretsMock.SecretsError)
		defer secretMock.DeactivateForceFailure()

		secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
			return "123", nil
		})

		err := restclient.AddMockedResponseFromFile("samples/jwks_response.json", jwtIssuer+"/auth/jwks", restclient.MethodGET)
		c.Nil(err)

		response, err := handleRequest(context.Background(), event)
		c.Equal(secretsMock.ErrForceFailure.Error(), response.Context["stringKey"])
		c.Contains(logMock.Output.String(), "getting_public_key_failed")
		logMock.Output.Reset()
	})

	t.Run("No tokens found for user", func(t *testing.T) {
		itMock := invalidtoken.InitDynamoMock()

		itMock.EmptyTable()

		err := restclient.AddMockedResponseFromFile("samples/jwks_response.json", jwtIssuer+"/auth/jwks", restclient.MethodGET)
		c.Nil(err)

		_, err = handleRequest(context.Background(), event)
		c.Nil(err)
		c.Contains(logMock.Output.String(), "no_tokens_found_for_user")
		logMock.Output.Reset()
	})

}

func dummyHandlerEvent() events.APIGatewayCustomAuthorizerRequest {
	return events.APIGatewayCustomAuthorizerRequest{
		Type:               "",
		AuthorizationToken: authToken,
		MethodArn:          "arn:aws:execute-api:us-east-1:811364018000:38qslpe8d9/ESTestInvoke-stage/GET/",
	}
}
