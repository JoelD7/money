package secrets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type SecretAPI interface {
	GetSecret(ctx context.Context, name string) (string, error)
}

type Secret struct{}

var (
	awsRegion = env.GetString("REGION", "us-east-1")

	ErrSecretNotFound = errors.New("secret not found")

	SecretClient SecretAPI

	secretCache *secretcache.Cache
)

func init() {
	SecretClient = &Secret{}

	sc, err := secretcache.New()
	if err != nil {
		panic(fmt.Errorf("secrets: %w", err))
	}

	secretCache = sc
}

func GetSecret(ctx context.Context, name string) (string, error) {
	return SecretClient.GetSecret(ctx, name)
}

func (s *Secret) GetSecret(ctx context.Context, name string) (string, error) {
	result, err := secretCache.GetSecretString(name)
	if err != nil && strings.Contains(err.Error(), "ResourceNotFoundException") {
		return "", ErrSecretNotFound
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

func CreateSecret(ctx context.Context, name, description string, secret []byte) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.CreateSecretInput{
		Name:                        aws.String(name),
		Description:                 aws.String(description),
		ForceOverwriteReplicaSecret: false,
		SecretBinary:                secret,
	}

	_, err = svc.CreateSecret(ctx, input)
	return err
}
