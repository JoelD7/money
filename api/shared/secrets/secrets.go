package secrets

import (
	"context"
	"errors"
	"github.com/JoelD7/money/api/shared/env"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretAPI interface {
	GetSecret(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error)
}

type Secret struct{}

var (
	awsRegion = env.GetString("REGION", "us-east-1")

	ErrSecretNotFound = errors.New("secret not found")

	SecretClient SecretAPI
)

func init() {
	SecretClient = &Secret{}
}

func GetSecret(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
	return SecretClient.GetSecret(ctx, name)
}

func (s *Secret) GetSecret(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ResourceNotFoundException") {
		return nil, ErrSecretNotFound
	}

	if err != nil {
		log.Fatal(err.Error())
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
