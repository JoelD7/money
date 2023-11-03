package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type createExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *createExpenseRequest) init() {
	dynamoClient := initDynamoClient()

	request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *createExpenseRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createExpenseHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(createExpenseRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *createExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	expense, err := validateInput(req, username)
	if err != nil {
		request.log.Error("validate_input_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	createExpense := usecases.NewExpenseCreator(request.expensesRepo, request.userRepo, request.periodRepo)

	newExpense, err := createExpense(ctx, username, expense)
	if err != nil {
		request.log.Error("create_expense_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusCreated, newExpense), nil
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

	return expense, nil
}
