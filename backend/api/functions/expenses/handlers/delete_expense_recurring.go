package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	expensesRecurring "github.com/JoelD7/money/backend/storage/expenses-recurring"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	derRequest *DeleteExpenseRecurringRequest
	derOnce    sync.Once
)

type DeleteExpenseRecurringRequest struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	Repo         expensesRecurring.Repository
}

func (request *DeleteExpenseRecurringRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error

	derOnce.Do(func() {
		request.Log = log
		dynamoClient := dynamo.InitClient(ctx)
		request.Repo, err = expensesRecurring.NewExpenseRecurringDynamoRepository(dynamoClient, envConfig.ExpensesRecurringTable)
	})

	request.startingTime = time.Now()

	return err
}

func (request *DeleteExpenseRecurringRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeleteExpenseRecurring(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if derRequest == nil {
		derRequest = new(DeleteExpenseRecurringRequest)
	}

	err := derRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("delete_expense_init_failed", err, req)
		return req.NewErrorResponse(err), nil
	}
	defer derRequest.finish()

	return derRequest.Process(ctx, req)
}

func (request *DeleteExpenseRecurringRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	expenseRecurringID, ok := req.PathParameters["expenseRecurringID"]
	if !ok || expenseRecurringID == "" {
		request.Log.Error("missing_expense_recurring_id", nil, req)

		return req.NewErrorResponse(models.ErrMissingExpenseRecurringID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.Log.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	deleteExpenseRecurring := usecases.NewExpenseRecurringEliminator(request.Repo)

	err = deleteExpenseRecurring(ctx, expenseRecurringID, username)
	if err != nil {
		request.Log.Error("delete_expense_recurring_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusNoContent, nil), nil
}
