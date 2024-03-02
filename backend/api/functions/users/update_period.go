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

type updatePeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *updatePeriodRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("update-period")
}

func (request *updatePeriodRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updatePeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updatePeriodRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updatePeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodBody, err := validateUpdateRequestBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	updatePeriod := usecases.NewPeriodUpdater(request.periodRepo)

	updatedPeriod, err := updatePeriod(ctx, periodBody.Username, periodBody.ID, periodBody)
	if err != nil {
		request.err = err
		request.log.Error("update_period_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, updatedPeriod), nil
}

func validateUpdateRequestBody(req *apigateway.Request) (*models.Period, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		return nil, err
	}

	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		return nil, models.ErrMissingPeriodID
	}

	periodModel := new(models.Period)

	err = json.Unmarshal([]byte(req.Body), periodModel)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrInvalidRequestBody, err)
	}

	periodModel.ID = periodID
	periodModel.Username = username

	if periodModel.Name == nil || periodModel.Name != nil && *periodModel.Name == "" {
		return nil, models.ErrMissingPeriodName
	}

	if periodModel.StartDate.IsZero() || periodModel.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	if periodModel.StartDate.After(periodModel.EndDate) {
		return nil, models.ErrStartDateShouldBeBeforeEndDate
	}

	if periodModel.CreatedDate.IsZero() {
		return nil, models.ErrMissingPeriodCreatedDate
	}

	return periodModel, nil
}
