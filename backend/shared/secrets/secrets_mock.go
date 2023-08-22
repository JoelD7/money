package secrets

import (
	"context"
	"errors"
)

type FailureCondition string

type MockSecret struct {
	secretResponders map[string]func(ctx context.Context, name string) (string, error)
	emulatingErrors  map[FailureCondition]error
	mockedErr        error
}

const (
	SecretsError FailureCondition = "secrets error"
)

var (
	errResponderNotRegistered = errors.New("mocks/secrets: responder not registered")
	errMockNotInitialized     = errors.New("mocks/secrets: mock not initialized")

	ErrForceFailure = errors.New("mocks/secrets: force failure")
)

func InitSecretMock() *MockSecret {
	mock := &MockSecret{
		secretResponders: make(map[string]func(ctx context.Context, name string) (string, error)),
		emulatingErrors: map[FailureCondition]error{
			SecretsError: ErrForceFailure,
		},
	}

	SecretClient = mock

	return mock
}

func (m *MockSecret) ActivateForceFailure(condition FailureCondition) {
	m.mockedErr = m.emulatingErrors[condition]
}

func (m *MockSecret) DeactivateForceFailure() {
	m.mockedErr = nil
}

func (m *MockSecret) GetSecret(ctx context.Context, name string) (string, error) {
	if m.mockedErr != nil {
		return "", m.mockedErr
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
