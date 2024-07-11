package env

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	secretName = "staging/money/env"
	// Only the /tmp directory is writable in AWS Lambda
	envFileName = "/tmp/.env"
)

func LoadEnv(ctx context.Context) error {
	f, err := os.Create(envFileName)
	if err != nil {
		return fmt.Errorf("cannot create .env file: %v", err)
	}

	sm := secrets.NewAWSSecretManager()

	val, err := sm.GetSecret(ctx, secretName)
	if err != nil {
		return fmt.Errorf("cannot get secret: %v", err)
	}

	_, err = f.WriteString(val)
	if err != nil {
		return fmt.Errorf("cannot write secrets to .env file: %v", err)
	}

	err = godotenv.Load(envFileName)
	if err != nil {
		return fmt.Errorf("cannot read enviroment configuration file: %v", err)
	}

	return nil
}

func LoadEnvTesting() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("cannot read enviroment configuration file: %v", err)
	}

	return nil
}

func GetString(varName, defaultValue string) string {
	val, _ := os.LookupEnv(varName)

	if val == "" {
		return defaultValue
	}

	return val
}

func GetInt(varName string, defaultValue int) int {
	val, _ := os.LookupEnv(varName)

	if val == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intValue
}
