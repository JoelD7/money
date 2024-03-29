package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type deleteExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
}

func (request *deleteExpenseRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
}

func (request *deleteExpenseRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func deleteExpenseHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(deleteExpenseRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *deleteExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	expenseID, ok := req.PathParameters["expenseID"]
	if !ok || expenseID == "" {
		request.log.Error("missing_expense_id", nil, []models.LoggerObject{req})

		return req.NewErrorResponse(models.ErrMissingExpenseID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	deleteExpense := usecases.NewExpensesDeleter(request.expensesRepo)

	err = deleteExpense(ctx, expenseID, username)
	if err != nil {
		request.log.Error("delete_expense_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
