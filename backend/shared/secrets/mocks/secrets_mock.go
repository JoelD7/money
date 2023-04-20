package mocks

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/shared/secrets"
)

type MockSecret struct {
	secretResponders map[string]func(ctx context.Context, name string) (string, error)
}

var (
	errResponderAlreadyRegistered = errors.New("mocks/secrets: responder is already registered")
	errResponderNotRegistered     = errors.New("mocks/secrets: responder not registered")
	errMockNotInitialized         = errors.New("mocks/secrets: mock not initialized")

	ErrForceFailure = errors.New("mocks/secrets: force failure")
)

var (
	ForceFailure bool
)

func InitSecretMock() *MockSecret {
	mock := &MockSecret{
		secretResponders: make(map[string]func(ctx context.Context, name string) (string, error)),
	}

	secrets.SecretClient = mock

	return mock
}

func (m *MockSecret) GetSecret(ctx context.Context, name string) (string, error) {
	if ForceFailure {
		return "", ErrForceFailure
	}

	responder, ok := m.secretResponders[name]
	if !ok {
		panic(errResponderNotRegistered)
	}

	return responder(ctx, name)
}

func (m *MockSecret) RegisterResponder(secretName string, responder func(ctx context.Context, name string) (string, error)) {
	if m == nil {
		panic(errMockNotInitialized)
	}

	m.secretResponders[secretName] = responder
}
