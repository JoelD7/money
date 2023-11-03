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

type getPeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *getPeriodRequest) init() {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getPeriodRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getPeriodHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getPeriodRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getPeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.err = models.ErrMissingPeriodID
		request.log.Error("missing_period_id", request.err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", request.err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	getPeriod := usecases.NewPeriodGetter(request.periodRepo)

	userPeriod, err := getPeriod(ctx, username, periodID)
	if err != nil {
		request.err = err
		request.log.Error("get_period_failed", request.err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, userPeriod), nil
}
