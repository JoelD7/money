package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/invalidtoken"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	authToken     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6InJlYWQgd3JpdGUiLCJpc3MiOiJodHRwczovLzM4cXNscGU4ZDkuZXhlY3V0ZS1hcGkudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vc3RhZ2luZyIsInN1YiI6InRlc3RAZ21haWwuY29tIiwiYXVkIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6MzAwMCIsImV4cCI6MTcwODI5OTA4OCwibmJmIjoxNjc3MTk2ODg4LCJpYXQiOjE2NzcxOTUwODh9.S_wnwVHTs_-T9zOkIFVIblfYYZ338kgUDclRi5nzgzLxzfqo_jrxKYXwLVeVkRNq1etO4B2RmyFPsLVHpC4cGS_Kr093eOzdWta0F8nj_hbTK2ZtuNP88X8oKaDadyCbXFw3M6dxm0la9kf20CZRxFsbtJ0MqPBqW9lp3B_XRz_pTAqMQnVbyfmbQBZiGBKpK5Ur1g043YAP5B2cd2C0ARGyyWw1UzXJBZbM_8KUFLUtndjZIn_uF3z8fLaH4hrnN3Gz_CnRIhgb6kbAWJ2OWsSJb4l15vgzdw2GvOWHU7MHqX6VoIwPVUzFTMDHkzfjDjhnKdWDj2bL-I-XXvZgSg"
	authTokenHash = "974df41534aaa82ed040cc75f0e5b4700094f79a5a168164288e95752ad43bf3"
)

func init() {
	restclient.NewMockRestClient()
}

func TestHandleRequest(t *testing.T) {
	c := require.New(t)

	mockRestClient := restclient.NewMockRestClient()
	redisRepository := cache.NewRepository(cache.NewRedisCacheMock())
	secretMock := secrets.NewSecretMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	req := &request{
		cacheRepo:      redisRepository,
		secretsManager: secretMock,
		client:         mockRestClient,
		log:            logMock,
	}

	event := dummyHandlerEvent()

	err := mockRestClient.AddMockedResponseFromFileNoUrl("samples/jwks_response.json", restclient.MethodGET)
	c.Nil(err)

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
		return "123", nil
	})

	response, err := req.process(ctx, event)
	c.Nil(err)
	c.Equal(Allow.String(), response.PolicyDocument.Statement[0].Effect)
}

func TestHandlerError(t *testing.T) {
	c := require.New(t)

	mockRestClient := restclient.NewMockRestClient()
	redisRepository := cache.NewRepository(cache.NewRedisCacheMock())
	secretMock := secrets.NewSecretMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	req := &request{
		cacheRepo:      redisRepository,
		secretsManager: secretMock,
		client:         mockRestClient,
		log:            logMock,
	}

	event := dummyHandlerEvent()
	ogToken := event.AuthorizationToken

	t.Run("Invalid token length", func(t *testing.T) {
		event.AuthorizationToken = "dummy"

		_, err := req.process(ctx, event)
		c.ErrorIs(err, models.ErrUnauthorized)
		c.Contains(logMock.Output.String(), "getting_token_payload_failed")
		logMock.Output.Reset()
	})

	t.Run("Payload decoding failed", func(t *testing.T) {
		event.AuthorizationToken = "Bearer dummy.dummy.token"

		_, err := req.process(ctx, event)
		c.ErrorIs(err, models.ErrUnauthorized)
		c.Contains(logMock.Output.String(), "getting_token_payload_failed")
		logMock.Output.Reset()
	})

	t.Run("Signing key not found", func(t *testing.T) {
		secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
			return "456", nil
		})

		err := mockRestClient.AddMockedResponseFromFileNoUrl("samples/jwks_response.json", restclient.MethodGET)
		c.Nil(err)

		event.AuthorizationToken = ogToken
		response, err := req.process(ctx, event)
		c.Nil(err)
		c.NotNil(response.Context["stringKey"])
		c.Equal(models.ErrSigningKeyNotFound.Error(), response.Context["stringKey"])
		c.Equal(Deny.String(), response.PolicyDocument.Statement[0].Effect)
		c.Contains(logMock.Output.String(), "getting_public_key_failed")
		logMock.Output.Reset()
	})

	t.Run("Getting public key failed", func(t *testing.T) {
		secretMock.ActivateForceFailure(secrets.SecretsError)
		defer secretMock.DeactivateForceFailure()

		secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
			return "123", nil
		})

		err := mockRestClient.AddMockedResponseFromFileNoUrl("samples/jwks_response.json", restclient.MethodGET)
		c.Nil(err)

		response, err := req.process(ctx, event)
		c.Nil(err)
		c.Equal(secrets.ErrForceFailure.Error(), response.Context["stringKey"])
		c.Equal(Deny.String(), response.PolicyDocument.Statement[0].Effect)
		c.Contains(logMock.Output.String(), "getting_public_key_failed")
		logMock.Output.Reset()
	})

	t.Run("Invalid token detected", func(t *testing.T) {
		err := mockRestClient.AddMockedResponseFromFileNoUrl("samples/jwks_response.json", restclient.MethodGET)
		c.Nil(err)

		err = redisRepository.AddInvalidToken(ctx, "test@gmail.com", authTokenHash, 0)
		c.Nil(err)

		response, err := req.process(ctx, event)
		c.Nil(err)
		c.Equal(Deny.String(), response.PolicyDocument.Statement[0].Effect)
		c.Contains(logMock.Output.String(), "invalid_token_use_detected")
		logMock.Output.Reset()
	})

	t.Run("JWT verification failed", func(t *testing.T) {
		_ = invalidtoken.InitDynamoMock()

		err := mockRestClient.AddMockedResponseFromFileNoUrl("samples/jwks_response.json", restclient.MethodGET)
		c.Nil(err)

		event.AuthorizationToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxNTE2MjM5MDIyLCJpc3MiOiJodHRwczovLzM4cXNscGU4ZDkuZXhlY3V0ZS1hcGkudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vc3RhZ2luZyJ9.uOmHNc9EwOQvu6qfeksVaDuqy4t8TmIGgoECUpPONnennzeDP-DgfH__kwwazENCRtjy75lbI7wbOdQjFL7qrcjopvF9NR4Ygf1S3nqPeCs4Db_i2XqD8KMzNEm8JxJ6iwJRZ26NrZEgrXIvJapBJ-JTaWKjKZdKYi5jjvVmrMNbvvDP-ZjUuOfFYrKWXZeyIhYT2YK3tdx48-dZn7JwWoGWZPAei99Fw-QzbGk9gaGOjv119-4JLVUfRDGOwibD4eGgoRQn3VZHgFwW-8cJod6XoQcmTuq_jHDRa28jwMIob6XGtMyMGqW5SNvhO6JigtmeaPY9jqLVdbXY_oGWbA"

		response, err := req.process(ctx, event)
		c.Nil(err)
		c.Equal(Deny.String(), response.PolicyDocument.Statement[0].Effect)
		c.Contains(logMock.Output.String(), "jwt_validation_failed")
		logMock.Output.Reset()
	})
}

func dummyHandlerEvent() events.APIGatewayCustomAuthorizerRequest {
	return events.APIGatewayCustomAuthorizerRequest{
		Type:               "",
		AuthorizationToken: "Bearer " + authToken,
		MethodArn:          "arn:aws:execute-api:us-east-1:811364018000:38qslpe8d9/ESTestInvoke-stage/GET/",
	}
}
