package main

import (
	"bytes"
	"context"
	secretsMock "github.com/JoelD7/money/api/shared/mocks/secrets"
	"github.com/JoelD7/money/api/shared/restclient"
	"github.com/JoelD7/money/auth/authenticator/secrets"
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
	restclient.Client = &restclient.MockClient{}
	secrets.SecretClient = &secretsMock.MockSecret{}
}

func TestHandleRequest(t *testing.T) {
	c := require.New(t)

	event := dummyHandlerEvent()

	data, err := os.ReadFile("samples/jwks_response.json")
	c.Nil(err)

	r := io.NopCloser(bytes.NewReader(data))

	restclient.GetFunction = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       r,
		}, nil
	}

	secretsMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
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
	event.AuthorizationToken = "dummy"

	_, err := handleRequest(context.Background(), event)
	c.ErrorIs(err, errUnauthorized)

	secretsMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String("123"),
		}, nil
	})
}

func dummyHandlerEvent() events.APIGatewayCustomAuthorizerRequest {
	return events.APIGatewayCustomAuthorizerRequest{
		Type:               "",
		AuthorizationToken: authToken,
		MethodArn:          "arn:aws:execute-api:us-east-1:811364018000:38qslpe8d9/ESTestInvoke-stage/GET/",
	}
}
