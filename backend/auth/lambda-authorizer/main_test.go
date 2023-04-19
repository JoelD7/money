package main

import (
	"bytes"
	"context"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	restMock "github.com/JoelD7/money/backend/shared/restclient/mocks"
	secretsMock "github.com/JoelD7/money/backend/shared/secrets/mocks"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

var (
	authToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6InJlYWQgd3JpdGUiLCJpc3MiOiJodHRwczovLzM4cXNscGU4ZDkuZXhlY3V0ZS1hcGkudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vc3RhZ2luZyIsInN1YiI6InRlc3RAZ21haWwuY29tIiwiYXVkIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6MzAwMCIsImV4cCI6MTcwODI5OTA4OCwibmJmIjoxNjc3MTk2ODg4LCJpYXQiOjE2NzcxOTUwODh9.S_wnwVHTs_-T9zOkIFVIblfYYZ338kgUDclRi5nzgzLxzfqo_jrxKYXwLVeVkRNq1etO4B2RmyFPsLVHpC4cGS_Kr093eOzdWta0F8nj_hbTK2ZtuNP88X8oKaDadyCbXFw3M6dxm0la9kf20CZRxFsbtJ0MqPBqW9lp3B_XRz_pTAqMQnVbyfmbQBZiGBKpK5Ur1g043YAP5B2cd2C0ARGyyWw1UzXJBZbM_8KUFLUtndjZIn_uF3z8fLaH4hrnN3Gz_CnRIhgb6kbAWJ2OWsSJb4l15vgzdw2GvOWHU7MHqX6VoIwPVUzFTMDHkzfjDjhnKdWDj2bL-I-XXvZgSg"
)

func init() {
	restclient.Client = &restMock.MockClient{}
	logger.InitLoggerMock()
}

func TestHandleRequest(t *testing.T) {
	c := require.New(t)

	event := dummyHandlerEvent()

	err := mockRestClientGetFromFile("samples/jwks_response.json")
	c.Nil(err)

	secretMock := secretsMock.InitSecretMock()

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("123"),
		}, nil
	})

	_, err = handleRequest(context.Background(), event)
	c.Nil(err)
}

func TestHandlerError(t *testing.T) {
	c := require.New(t)

	event := dummyHandlerEvent()
	authToken := event.AuthorizationToken

	event.AuthorizationToken = "dummy"

	_, err := handleRequest(context.Background(), event)
	c.ErrorIs(err, errUnauthorized)

	event.AuthorizationToken = "Bearer dummy.dummy.token"
	_, err = handleRequest(context.Background(), event)
	c.ErrorIs(err, errUnauthorized)

	secretMock := secretsMock.InitSecretMock()

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("456"),
		}, nil
	})

	err = mockRestClientGetFromFile("samples/jwks_response.json")
	c.Nil(err)

	event.AuthorizationToken = authToken
	_, err = handleRequest(context.Background(), event)
	c.ErrorIs(err, errSigningKeyNotFound)

	secretsMock.ForceFailure = true
	defer func() {
		secretsMock.ForceFailure = false
	}()

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("123"),
		}, nil
	})

	err = mockRestClientGetFromFile("samples/jwks_response.json")
	c.Nil(err)

	_, err = handleRequest(context.Background(), event)
	c.ErrorIs(err, secretsMock.ErrForceFailure)
}

func dummyHandlerEvent() events.APIGatewayCustomAuthorizerRequest {
	return events.APIGatewayCustomAuthorizerRequest{
		Type:               "",
		AuthorizationToken: authToken,
		MethodArn:          "arn:aws:execute-api:us-east-1:811364018000:38qslpe8d9/ESTestInvoke-stage/GET/",
	}
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
