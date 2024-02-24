package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type deletePeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *deletePeriodRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("delete-period")
}

func (request *deletePeriodRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func deletePeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(deletePeriodRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *deletePeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.log.Error("missing_period_id", nil, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	deletePeriod := usecases.NewPeriodDeleter(request.periodRepo)

	err = deletePeriod(ctx, periodID, username)
	if err != nil {
		request.log.Error("delete_period_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
