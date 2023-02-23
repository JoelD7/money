package secrets

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type MockSecret struct{}

var (
	//GetSecretFunc    func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error)
	secretResponders = make(map[string]func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error))

	errResponderAlreadyRegistered = errors.New("secrets: responder is already registered")
	errResponderNotRegistered     = errors.New("secrets: responder not registered")
)

func (m *MockSecret) GetSecret(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
	responder, ok := secretResponders[name]
	if !ok {
		panic(errResponderNotRegistered)
	}

	return responder(ctx, name)
}

func RegisterResponder(secretName string, responder func(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error)) {
	if _, ok := secretResponders[secretName]; ok {
		panic(errResponderAlreadyRegistered)
	}

	secretResponders[secretName] = responder
}
