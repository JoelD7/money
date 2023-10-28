package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"time"
)

type updatePeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *updatePeriodRequest) init() {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *updatePeriodRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updatePeriodHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updatePeriodRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updatePeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	_, err := validateUpdateRequestBody(req)
	if err != nil {
		request.err = err
		request.log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return nil, nil
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

	if periodModel.Name == nil {
		return nil, models.ErrMissingPeriodName
	}

	if periodModel.StartDate.IsZero() || periodModel.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	if periodModel.StartDate.After(periodModel.EndDate.Time) {
		return nil, models.ErrStartDateShouldBeBeforeEndDate
	}

	if periodModel.CreatedDate.IsZero() {
		return nil, models.ErrMissingPeriodCreatedDate
	}

	if periodModel.UpdatedDate.IsZero() {
		return nil, models.ErrMissingPeriodUpdatedDate
	}

	return nil, nil
}
