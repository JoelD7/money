package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	ceRequest *createExpenseRequest
	ceOnce    sync.Once
)

type createExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *createExpenseRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	ceOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)
		request.log = log

		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig.ExpensesTable, envConfig.ExpensesRecurringTable, envConfig.PeriodUserExpenseIndex)
		if err != nil {
			return
		}

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *createExpenseRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateExpense(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if ceRequest == nil {
		ceRequest = new(createExpenseRequest)
	}

	err := ceRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("create_expense_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	defer ceRequest.finish()

	return ceRequest.process(ctx, req)
}

func (request *createExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	expense, err := validateInput(req, username)
	if err != nil {
		request.log.Error("validate_input_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createExpense := usecases.NewExpenseCreator(request.expensesRepo, request.periodRepo)

	newExpense, err := createExpense(ctx, username, expense)
	if err != nil {
		request.log.Error("create_expense_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, newExpense), nil
}

func validateInput(req *apigateway.Request, username string) (*models.Expense, error) {
	expense := new(models.Expense)

	err := json.Unmarshal([]byte(req.Body), expense)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	err = validate.Email(username)
	if err != nil {
		return nil, err
	}

	expense.Username = username

	if expense.Name == nil {
		return nil, models.ErrMissingName
	}

	if expense.Amount == nil {
		return nil, models.ErrMissingAmount
	}

	err = validate.Amount(expense.Amount)
	if err != nil {
		return nil, err
	}

	if expense.Period == "" {
		return nil, models.ErrMissingPeriod
	}

	if expense.IsRecurring && expense.RecurringDay == nil {
		return nil, models.ErrMissingRecurringDay
	}

	if (expense.IsRecurring && expense.RecurringDay != nil) && (*expense.RecurringDay < 1 || *expense.RecurringDay > 31) {
		return nil, models.ErrInvalidRecurringDay
	}

	return expense, nil
}
