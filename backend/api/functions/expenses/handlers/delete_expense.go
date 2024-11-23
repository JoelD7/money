package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var deRequest *deleteExpenseRequest
var deOnce sync.Once

type deleteExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
}

func (request *deleteExpenseRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error

	deOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.log = log
		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
	})

	request.startingTime = time.Now()

	return err
}

func (request *deleteExpenseRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeleteExpense(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if deRequest == nil {
		deRequest = new(deleteExpenseRequest)
	}

	err := deRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("delete_expense_init_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}
	defer deRequest.finish()

	return deRequest.process(ctx, req)
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
