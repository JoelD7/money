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
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *updateExpenseRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	ueOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *updateExpenseRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func UpdateExpense(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if ueRequest == nil {
		ueRequest = new(updateExpenseRequest)
	}

	err := ueRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("update_expense_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer ueRequest.finish()

	return ueRequest.process(ctx, req)
}

func (request *updateExpenseRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	expenseID, ok := req.PathParameters["expenseID"]
	if !ok || expenseID == "" {
		logger.Error("missing_expense_id", nil, req)

		return req.NewErrorResponse(models.ErrMissingExpenseID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	expense, err := validateUpdateInput(req, username)
	if err != nil {
		logger.Error("validate_input_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	updateExpense := usecases.NewExpenseUpdater(request.expensesRepo, request.periodRepo, request.userRepo)

	updatedExpense, err := updateExpense(ctx, expenseID, username, expense)
	if err != nil {
		logger.Error("update_expense_failed", err, req)

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
