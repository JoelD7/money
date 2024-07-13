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

var ueRequest *updateExpenseRequest
var ueOnce sync.Once

type updateExpenseRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *updateExpenseRequest) init(ctx context.Context, log logger.LogAPI) {
	ueOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
		request.periodRepo = period.NewDynamoRepository(dynamoClient)
		request.userRepo = users.NewDynamoRepository(dynamoClient)
		request.log = log
	})
	request.startingTime = time.Now()
}

func (request *updateExpenseRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func UpdateExpense(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if ueRequest == nil {
		ueRequest = new(updateExpenseRequest)
	}

	ueRequest.init(ctx, log)
	defer ueRequest.finish()

	return ueRequest.process(ctx, req)
}

func (request *updateExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
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

	expense, err := validateUpdateInput(req, username)
	if err != nil {
		request.log.Error("validate_input_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	updateExpense := usecases.NewExpenseUpdater(request.expensesRepo, request.periodRepo, request.userRepo)

	updatedExpense, err := updateExpense(ctx, expenseID, username, expense)
	if err != nil {
		request.log.Error("update_expense_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, updatedExpense), nil
}

func validateUpdateInput(req *apigateway.Request, username string) (*models.Expense, error) {
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

	err = validate.Amount(expense.Amount)
	if err != nil {
		return nil, err
	}

	return expense, nil
}
