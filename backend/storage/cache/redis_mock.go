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

func (r *redisMock) GetInvalidTokens(ctx context.Context, username string) ([]*models.InvalidToken, error) {
	if r.mockedErr != nil {
		return nil, r.mockedErr
	}

	invalidTokens, ok := r.store[username]
	if !ok {
		return nil, models.ErrInvalidTokensNotFound
	}

	return invalidTokens, nil
}

func (r *redisMock) AddInvalidToken(ctx context.Context, username, token string, ttl int64) error {
	if r.mockedErr != nil {
		return r.mockedErr
	}

	invalidTokens, ok := r.store[username]
	if !ok {
		r.store[username] = append(r.store[username], &models.InvalidToken{Token: token})
		return nil
	}

	r.store[username] = append(invalidTokens, &models.InvalidToken{Token: token})

	return nil
}

func (r *redisMock) DeleteInvalidToken(username string) {
	delete(r.store, username)
}

func (r *redisMock) AddResource(ctx context.Context, key string, resource interface{}, ttl int64) error {
	return nil
}

func (r *redisMock) GetResource(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (r *redisMock) SetTTL(ttl int64) {
	return
}

func (r *redisMock) AddIncomePeriods(ctx context.Context, username string, periods []string) error {
	return nil
}

func (r *redisMock) GetIncomePeriods(ctx context.Context, username string) ([]string, error) {
	return nil, nil
}

func (r *redisMock) DeleteIncomePeriods(ctx context.Context, username string, periods ...string) error {
	return nil
}
