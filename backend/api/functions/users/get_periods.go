package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type getPeriodsResponse struct {
	Periods []*models.Period `json:"periods"`
	NextKey string           `json:"next_key"`
}

type getPeriodsRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
	username     string
	startKey     string
	pageSize     int
}

func (request *getPeriodsRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("get-periods")
}

func (request *getPeriodsRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getPeriodsHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getPeriodsRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getPeriodsRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	err := request.prepareRequest(req)
	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	getPeriods := usecases.NewPeriodsGetter(request.periodRepo)

	userPeriods, nextKey, err := getPeriods(ctx, request.username, request.startKey, request.pageSize)
	if err != nil {
		request.err = err
		request.log.Error("get_periods_failed", request.err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	res := &getPeriodsResponse{
		Periods: userPeriods,
		NextKey: nextKey,
	}

	return apigateway.NewJSONResponse(http.StatusOK, res), nil
}

func (request *getPeriodsRequest) prepareRequest(req *apigateway.Request) error {
	var err error

	request.username, err = apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return err
	}

	err = validate.Email(request.username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{
			request.log.MapToLoggerObject("user_data", map[string]interface{}{
				"s_username": request.username,
			}),
		})

		return err
	}

	request.startKey, request.pageSize, err = getRequestQueryParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return err
	}

	return nil
}
