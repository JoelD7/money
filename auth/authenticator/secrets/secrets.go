package secrets

import (
	"context"
	"github.com/JoelD7/money/api/shared/env"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	awsRegion = env.GetString("REGION", "us-east-1")
)

func init() {

}

func GetSecret(ctx context.Context, name string) (*secretsmanager.GetSecretValueOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		log.Fatal(err.Error())
	}

	return result, nil

}
