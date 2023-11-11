package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type createIncomeRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	incomeRepo   income.Repository
	periodRepo   period.Repository
}

func (request *createIncomeRequest) init() {
	dynamoClient := initDynamoClient()

	request.incomeRepo = income.NewDynamoRepository(dynamoClient)
	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *createIncomeRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createIncomeHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(createIncomeRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *createIncomeRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	reqIncome, err := validateCreateIncomeBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_create_income_body_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	createIncome := usecases.NewIncomeCreator(request.incomeRepo, request.periodRepo)

	newIncome, err := createIncome(ctx, username, reqIncome)
	if err != nil {
		request.err = err
		request.log.Error("create_income_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, newIncome), nil
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