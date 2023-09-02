package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type SecretManager interface {
	GetSecret(ctx context.Context, name string) (string, error)
}

type AWSSecretManager struct {
	secretCache *secretcache.Cache
}

var (
	secretCache *secretcache.Cache
)

func init() {
	sc, err := secretcache.New()
	if err != nil {
		panic(fmt.Errorf("secrets: %w", err))
	}

	secretCache = sc
}

func NewAWSSecretManager() *AWSSecretManager {
	return &AWSSecretManager{secretCache: secretCache}
}

func (s *AWSSecretManager) GetSecret(ctx context.Context, name string) (string, error) {
	result, err := secretCache.GetSecretString(name)
	if err != nil && strings.Contains(err.Error(), "ResourceNotFoundException") {
		return "", models.ErrSecretNotFound
	}

	if err != nil {
		return "", err
	}

	return result, nil
}
