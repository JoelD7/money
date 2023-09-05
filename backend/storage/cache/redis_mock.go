package cache

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type redisMock struct {
	store     map[string][]*models.InvalidToken
	mockedErr error
}

// NewRedisCacheMock creates a redis mock by mocking the underlying redis client.
func NewRedisCacheMock() *redisMock {
	return &redisMock{
		store: make(map[string][]*models.InvalidToken),
	}
}

func (r *redisMock) ActivateForceFailure(err error) {
	r.mockedErr = err
}

func (r *redisMock) DeactivateForceFailure() {
	r.mockedErr = nil
}

func (r *redisMock) GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error) {
	if r.mockedErr != nil {
		return nil, r.mockedErr
	}

	invalidTokens, ok := r.store[email]
	if !ok {
		return nil, models.ErrInvalidTokensNotFound
	}

	return invalidTokens, nil
}

func (r *redisMock) AddInvalidToken(ctx context.Context, email, token string, ttl int64) error {
	if r.mockedErr != nil {
		return r.mockedErr
	}

	invalidTokens, ok := r.store[email]
	if !ok {
		r.store[email] = append(r.store[email], &models.InvalidToken{Token: token})
		return nil
	}

	r.store[email] = append(invalidTokens, &models.InvalidToken{Token: token})

	return nil
}

func (r *redisMock) DeleteInvalidToken(email string) {
	delete(r.store, email)
}
