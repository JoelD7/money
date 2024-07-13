package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	errMissingExpenseID = apigateway.NewError("missing expense ID", http.StatusBadRequest)
	geExpenseRequest    *getExpenseRequest
	getExpenseOnce      sync.Once
)

type getExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
}

func (request *getExpenseRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	getExpenseOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, tableName, expensesRecurringTableName)
		if err != nil {
			return
		}
		request.userRepo = users.NewDynamoRepository(dynamoClient)
		request.log = log
	})
	request.startingTime = time.Now()

	return err
}

func (request *getExpenseRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetExpense(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if geExpenseRequest == nil {
		geExpenseRequest = new(getExpenseRequest)
	}

	err := geExpenseRequest.init(ctx, log)
	if err != nil {
		log.Error("get_expense_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	defer geExpenseRequest.finish()

	return geExpenseRequest.process(ctx, req)
}

func (request *getExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	expenseID, ok := req.PathParameters["expenseID"]
	if !ok || expenseID == "" {
		request.log.Error("missing_expense_id", errMissingExpenseID, []models.LoggerObject{req})

		return req.NewErrorResponse(errMissingExpenseID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{
			request.log.MapToLoggerObject("user_data", map[string]interface{}{
				"s_username": username,
			}),
		})

		return req.NewErrorResponse(err), nil
	}

	getExpense := usecases.NewExpenseGetter(request.expensesRepo, request.userRepo)

	expense, err := getExpense(ctx, username, expenseID)
	if err != nil {
		request.log.Error("get_expense_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, expense), nil
}
