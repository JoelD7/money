package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var preRequest *PatchRecurrentExpenseRequest
var preOnce sync.Once

type PatchRecurrentExpenseRequest struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	ExpensesRepo expenses.Repository
	PeriodRepo   period.Repository
}

type patchRecurrentExpenseRequestBody struct {
	Period string `json:"period"`
}

func (request *PatchRecurrentExpenseRequest) init(log logger.LogAPI) {
	preOnce.Do(func() {
		dynamoClient := initDynamoClient()

		request.ExpensesRepo = expenses.NewDynamoRepository(dynamoClient)
		request.PeriodRepo = period.NewDynamoRepository(dynamoClient)
		request.Log = log
	})
	request.startingTime = time.Now()
}

func (request *PatchRecurrentExpenseRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func PatchRecurrentExpense(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if preRequest == nil {
		preRequest = new(PatchRecurrentExpenseRequest)
	}

	preRequest.init(log)
	defer preRequest.finish()

	return preRequest.Process(ctx, req)
}

func (request *PatchRecurrentExpenseRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.Log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	reqBody, err := validateRequestBody(req)
	if err != nil {
		request.Log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	updateExpensesWoPeriod := usecases.NewExpensesPeriodSetter(request.ExpensesRepo, request.PeriodRepo, request.Log)

	err = updateExpensesWoPeriod(ctx, username, reqBody.Period)
	if err != nil {
		request.Log.Error("patch_recurrent_expenses_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, ""), nil
}

func validateRequestBody(req *apigateway.Request) (*patchRecurrentExpenseRequestBody, error) {
	requestBody := new(patchRecurrentExpenseRequestBody)

	err := json.Unmarshal([]byte(req.Body), requestBody)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, models.ErrInvalidRequestBody)
	}

	if requestBody.Period == "" {
		return nil, models.ErrMissingPeriod
	}

	return requestBody, nil
}
