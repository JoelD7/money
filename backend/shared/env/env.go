package env

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	secretName = "staging/money/env"
	// Only the /tmp directory is writable in AWS Lambda
	envFileName = "/tmp/.env"
)

func LoadEnv(ctx context.Context) (*models.EnvironmentConfiguration, error) {
	f, err := os.Create(envFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot create .env file: %v", err)
	}

	sm := secrets.NewAWSSecretManager()

	secret, err := sm.GetSecret(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("cannot get secret: %v", err)
	}

	_, err = f.WriteString(secret)
	if err != nil {
		return nil, fmt.Errorf("cannot write secrets to .env file: %v", err)
	}

	err = godotenv.Load(envFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read enviroment configuration file: %v", err)
	}

	return getEnvironmentConfig(), nil
}

func getEnvironmentConfig() *models.EnvironmentConfiguration {
	return &models.EnvironmentConfiguration{
		MissingExpensePeriodQueueURL: GetString("MISSING_EXPENSE_PERIOD_QUEUE_URL", ""),
		AwsRegion:                    GetString("AWS_REGION", ""),
		LogstashType:                 GetString("LOGSTASH_TYPE", ""),
		LogstashHost:                 GetString("LOGSTASH_HOST", ""),
		LogstashPort:                 GetString("LOGSTASH_PORT", ""),
		RedisURL:                     GetString("REDIS_URL", ""),
		CorsOrigin:                   GetString("CORS_ORIGIN", ""),
		AccessTokenDuration:          GetString("ACCESS_TOKEN_DURATION", ""),
		RefreshTokenDuration:         GetString("REFRESH_TOKEN_DURATION", ""),
		TokenAudience:                GetString("TOKEN_AUDIENCE", ""),
		TokenIssuer:                  GetString("TOKEN_ISSUER", ""),
		TokenPrivateSecret:           GetString("TOKEN_PRIVATE_SECRET", ""),
		TokenPublicSecret:            GetString("TOKEN_PUBLIC_SECRET", ""),
		KidSecret:                    GetString("KID_SECRET", ""),
		TokenScope:                   GetString("TOKEN_SCOPE", ""),
		LambdaTimeout:                GetString("LAMBDA_TIMEOUT", ""),
		UsersTable:                   GetString("USERS_TABLE_NAME", ""),
		ExpensesTable:                GetString("EXPENSES_TABLE_NAME", ""),
		ExpensesRecurringTable:       GetString("EXPENSES_RECURRING_TABLE_NAME", ""),
		IncomeTable:                  GetString("INCOME_TABLE_NAME", ""),
		InvalidTokenTable:            GetString("INVALID_TOKEN_TABLE_NAME", ""),
		PeriodTable:                  GetString("PERIOD_TABLE_NAME", ""),
		UniquePeriodTable:            GetString("UNIQUE_PERIOD_TABLE_NAME", ""),
		SavingsTable:                 GetString("SAVINGS_TABLE_NAME", ""),
		PeriodSavingIndexName:        GetString("PERIOD_SAVING_INDEX_NAME", ""),
		SavingGoalSavingIndexName:    GetString("SAVING_GOAL_SAVING_INDEX_NAME", ""),
		SavingGoalsTable:             GetString("SAVING_GOALS_TABLE_NAME", ""),
		BatchWriteRetries:            GetInt("BATCH_WRITE_RETRIES", 3),
		BatchWriteBaseDelayInMs:      GetInt("BATCH_WRITE_BASE_DELAY_IN_MS", 300),
		BatchWriteBackoffFactor:      GetInt("BATCH_WRITE_BACKOFF_FACTOR", 2),
		DynamodbMaxBatchWrite:        GetInt("DYNAMODB_MAX_BATCH_WRITE", 25),
		PeriodUserCreatedDateIndex:   GetString("PERIOD_USER_CREATED_DATE_INDEX", ""),
		PeriodUserExpenseIndex:       GetString("PERIOD_USER_EXPENSE_INDEX", ""),
		UsernameCreatedDateIndex:     GetString("USERNAME_CREATED_DATE_INDEX", ""),
		PeriodUserIncomeIndex:        GetString("PERIOD_USER_INCOME_INDEX", ""),
		PeriodUserNameExpenseIDIndex: GetString("PERIOD_USER_NAME_EXPENSE_ID_INDEX", ""),
		PeriodUserAmountIndex:        GetString("PERIOD_USER_AMOUNT_INDEX", ""),
		PeriodUserNameIncomeIDIndex:  GetString("PERIOD_USER_NAME_INCOME_ID_INDEX", ""),
		UsernameTargetIndex:          GetString("USERNAME_TARGET_INDEX", ""),
		UsernameDeadlineIndex:        GetString("USERNAME_DEADLINE_INDEX", ""),
		UsernameAmountIndex:          GetString("USERNAME_AMOUNT_INDEX", ""),
	}
}

// LoadEnvTesting loads the environment variables from the .env file for testing purposes.
func LoadEnvTesting() error {
	//Currently, it appears that godotenv doesn't support loading files using relative paths. This is why we need to use
	//absolute paths to load the .env file.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("unable to identify current directory needed to load .env")
	}

	basepath := filepath.Dir(file)
	path := filepath.Join(basepath, "../../.env")
	err := godotenv.Load(path)
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

func GetBool(varName string) bool {
	val, _ := os.LookupEnv(varName)

	return val == "true"
}
