package models

type EnvironmentConfiguration struct {
	MissingExpensePeriodQueueURL string `json:"MISSING_EXPENSE_PERIOD_QUEUE_URL"`
	AwsRegion                    string `json:"AWS_REGION"`

	LogstashType string `json:"LOGSTASH_TYPE"`
	LogstashHost string `json:"LOGSTASH_HOST"`
	LogstashPort string `json:"LOGSTASH_PORT"`

	RedisURL   string `json:"REDIS_URL"`
	CorsOrigin string `json:"CORS_ORIGIN"`

	AccessTokenDuration  string `json:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration string `json:"REFRESH_TOKEN_DURATION"`
	TokenAudience        string `json:"TOKEN_AUDIENCE"`
	TokenIssuer          string `json:"TOKEN_ISSUER"`
	TokenPrivateSecret   string `json:"TOKEN_PRIVATE_SECRET"`
	TokenPublicSecret    string `json:"TOKEN_PUBLIC_SECRET"`
	KidSecret            string `json:"KID_SECRET"`
	TokenScope           string `json:"TOKEN_SCOPE"`
	LambdaTimeout        string `json:"LAMBDA_TIMEOUT"`

	UsersTable                   string `json:"USERS_TABLE_NAME"`
	ExpensesTable                string `json:"EXPENSES_TABLE_NAME"`
	ExpensesRecurringTable       string `json:"EXPENSES_RECURRING_TABLE_NAME"`
	IncomeTable                  string `json:"INCOME_TABLE_NAME"`
	PeriodUserIncomeIndex        string `json:"PERIOD_USER_INCOME_INDEX"`
	InvalidTokenTable            string `json:"INVALID_TOKEN_TABLE_NAME"`
	PeriodTable                  string `json:"PERIOD_TABLE_NAME"`
	UniquePeriodTable            string `json:"UNIQUE_PERIOD_TABLE_NAME"`
	SavingsTable                 string `json:"SAVINGS_TABLE_NAME"`
	PeriodSavingIndexName        string `json:"PERIOD_SAVING_INDEX_NAME"`
	PeriodUserExpenseIndex       string `json:"PERIOD_USER_EXPENSE_INDEX"`
	SavingGoalSavingIndexName    string `json:"SAVING_GOAL_SAVING_INDEX_NAME"`
	SavingGoalsTable             string `json:"SAVING_GOALS_TABLE_NAME"`
	PeriodUserCreatedDateIndex   string `json:"PERIOD_USER_CREATED_DATE_INDEX"`
	UsernameCreatedDateIndex     string `json:"USERNAME_CREATED_DATE_INDEX"`
	PeriodUserNameExpenseIDIndex string `json:"USERNAME_NAME_EXPENSE_ID_INDEX"`
	PeriodUserAmountIndex        string `json:"USERNAME_AMOUNT_INDEX"`

	BatchWriteRetries       int `json:"BATCH_WRITE_RETRIES"`
	BatchWriteBaseDelayInMs int `json:"BATCH_WRITE_BASE_DELAY_IN_MS"`
	BatchWriteBackoffFactor int `json:"BATCH_WRITE_BACKOFF_FACTOR"`
	DynamodbMaxBatchWrite   int `json:"DYNAMODB_MAX_BATCH_WRITE"`
}
