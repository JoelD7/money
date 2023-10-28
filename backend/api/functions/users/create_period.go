package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type createPeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *createPeriodRequest) init() {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *createPeriodRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createPeriodHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(createPeriodRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *createPeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodModel, err := request.validateCreateRequestBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	createPeriod := usecases.NewPeriodCreator(request.periodRepo)

	createdPeriod, err := createPeriod(ctx, username, periodModel)
	if err != nil {
		request.err = err
		request.log.Error("create_period_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusCreated, createdPeriod), nil
}

func (request *createPeriodRequest) validateCreateRequestBody(req *apigateway.Request) (*models.Period, error) {
	p := new(models.Period)

	err := json.Unmarshal([]byte(req.Body), p)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, models.ErrInvalidRequestBody)
	}

	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	return p, nil
}
