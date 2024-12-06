package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"github.com/aws/aws-lambda-go/events"
	"sync"
	"time"
)

var (
	preRequest *Request
	preOnce    sync.Once
)

type Request struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	ExpensesRepo expenses.Repository
	PeriodRepo   period.Repository
}

func (request *Request) init(ctx context.Context) error {
	var err error

	preOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)
		request.Log = logger.NewLogger()

		periodTableNameEnv := env.GetString("PERIOD_TABLE_NAME", "")
		uniquePeriodTableNameEnv := env.GetString("UNIQUE_PERIOD_TABLE_NAME", "")

		envConfig := &models.EnvironmentConfiguration{
			ExpensesTable:          env.GetString("EXPENSES_TABLE_NAME", ""),
			ExpensesRecurringTable: env.GetString("EXPENSES_RECURRING_TABLE_NAME", ""),
			PeriodTable:            env.GetString("PERIOD_TABLE_NAME", ""),
			UniquePeriodTable:      env.GetString("UNIQUE_PERIOD_TABLE_NAME", ""),
			PeriodUserExpenseIndex: env.GetString("PERIOD_USER_EXPENSE_INDEX", ""),
		}

		request.ExpensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.PeriodRepo, err = period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *Request) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	if preRequest == nil {
		preRequest = &Request{}
	}

	err := preRequest.init(ctx)
	if err != nil {
		preRequest.Log.Error("init_failed", err, nil)

		return err
	}
	defer preRequest.finish()

	for _, record := range sqsEvent.Records {
		err = preRequest.ProcessMessage(ctx, models.SQSMessage{SQSMessage: record})
		if err != nil {
			return err
		}
	}

	preRequest.Log.Info("message_processing_successful", []models.LoggerObject{
		preRequest.Log.MapToLoggerObject("message_data", map[string]interface{}{
			"i_message_count": len(sqsEvent.Records),
		}),
	})

	return nil
}

func (request *Request) ProcessMessage(ctx context.Context, record models.SQSMessage) error {
	msgBody, err := validateMessageBody(record)
	if err != nil {
		request.Log.Error("validate_request_body_failed", err, []models.LoggerObject{record})

		return err
	}

	updateExpensesWoPeriod := usecases.NewExpensesPeriodSetter(request.ExpensesRepo, request.PeriodRepo, request.Log)

	err = updateExpensesWoPeriod(ctx, msgBody.Username, msgBody.Period)
	if err != nil {
		request.Log.Error("patch_recurrent_expenses_failed", err, []models.LoggerObject{record})

		return err
	}

	return nil
}

func validateMessageBody(record models.SQSMessage) (*models.MissingExpensePeriodMessage, error) {
	msgBody := new(models.MissingExpensePeriodMessage)

	err := json.Unmarshal([]byte(record.Body), msgBody)
	if err != nil {
		return nil, fmt.Errorf("invalid message body: %v", err)
	}

	if msgBody.Username == "" {
		return nil, models.ErrMissingUsername
	}

	if msgBody.Period == "" {
		return nil, models.ErrMissingPeriod
	}

	return msgBody, nil
}
