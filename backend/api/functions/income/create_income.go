package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	tableName                = env.GetString("INCOME_TABLE_NAME", "")
	periodTableNameEnv       = env.GetString("PERIOD_TABLE_NAME", "")
	uniquePeriodTableNameEnv = env.GetString("UNIQUE_PERIOD_TABLE_NAME", "")

	ciRequest *createIncomeRequest
	ciOnce    sync.Once
)

type createIncomeRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	incomeRepo   income.Repository
	periodRepo   period.Repository
}

func (request *createIncomeRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	ciOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)
		request.log = log

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, tableName)
		if err != nil {
			return
		}
		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *createIncomeRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createIncomeHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if ciRequest == nil {
		ciRequest = new(createIncomeRequest)
	}

	ciRequest.init(ctx, log)
	defer ciRequest.finish()

	return ciRequest.process(ctx, req)
}

func (request *createIncomeRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	reqIncome, err := validateCreateIncomeBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_create_income_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createIncome := usecases.NewIncomeCreator(request.incomeRepo, request.periodRepo)

	newIncome, err := createIncome(ctx, username, reqIncome)
	if err != nil {
		request.err = err
		request.log.Error("create_income_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, newIncome), nil
}

func validateCreateIncomeBody(req *apigateway.Request) (*models.Income, error) {
	reqIncome := new(models.Income)

	err := json.Unmarshal([]byte(req.Body), &reqIncome)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	if reqIncome.Amount == nil {
		return nil, models.ErrMissingAmount
	}

	if reqIncome.Name == nil {
		return nil, models.ErrMissingName
	}

	if reqIncome.Period == nil {
		return nil, models.ErrMissingPeriod
	}

	err = validate.Amount(reqIncome.Amount)
	if err != nil {
		return nil, err
	}

	return reqIncome, nil
}
