package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type updateExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
}

func (request *updateExpenseRequest) init() {
	dynamoClient := initDynamoClient()

	request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *updateExpenseRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updateExpenseHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updateExpenseRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updateExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	expenseID, ok := req.PathParameters["expenseID"]
	if !ok || expenseID == "" {
		request.log.Error("missing_expense_id", nil, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(models.ErrMissingExpenseID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	expense, err := validateUpdateInput(req, username)
	if err != nil {
		request.log.Error("validate_input_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	updateExpense := usecases.NewExpenseUpdater(request.expensesRepo)

	err = updateExpense(ctx, expenseID, username, expense)
	if err != nil {
		request.log.Error("update_expense_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func validateUpdateInput(req *apigateway.Request, username string) (*models.Expense, error) {
	expense := new(models.Expense)

	err := json.Unmarshal([]byte(req.Body), expense)
	if err != nil {
		return nil, err
	}

	err = validate.Email(username)
	if err != nil {
		return nil, err
	}

	expense.Username = username

	err = validate.Amount(expense.Amount)
	if err != nil {
		return nil, err
	}

	return expense, nil
}
