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
	"sync"
	"time"
)

var cpRequest *createPeriodRequest
var cpOnce sync.Once

type createPeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *createPeriodRequest) init(log logger.LogAPI) {
	cpOnce.Do(func() {
		dynamoClient := initDynamoClient()

		request.periodRepo = period.NewDynamoRepository(dynamoClient)
		request.log = log
		request.log.SetHandler("create-period")
	})
	request.startingTime = time.Now()
}

func (request *createPeriodRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createPeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if cpRequest == nil {
		cpRequest = new(createPeriodRequest)
	}

	cpRequest.init(log)
	defer cpRequest.finish()

	return cpRequest.process(ctx, req)
}

func (request *createPeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodModel, err := request.validateCreateRequestBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createPeriod := usecases.NewPeriodCreator(request.periodRepo, request.log)

	createdPeriod, err := createPeriod(ctx, username, periodModel)
	if err != nil {
		request.err = err
		request.log.Error("create_period_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, createdPeriod), nil
}

func (request *createPeriodRequest) validateCreateRequestBody(req *apigateway.Request) (*models.Period, error) {
	p := new(models.Period)

	err := json.Unmarshal([]byte(req.Body), p)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, models.ErrInvalidRequestBody)
	}

	if p.Name == nil || p.Name != nil && *p.Name == "" {
		return nil, models.ErrMissingPeriodName
	}

	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	return p, nil
}
