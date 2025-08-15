package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"strings"
	"sync"
	"time"
)

var gpsRequest *getPeriodsRequest
var gpsOnce sync.Once

type getPeriodsResponse struct {
	Periods []*models.Period `json:"periods"`
	NextKey string           `json:"next_key"`
}

type getPeriodsRequest struct {
	startingTime time.Time
	err          error
	username     string
	Active       bool

	log        logger.LogAPI
	periodRepo period.Repository
	*models.QueryParameters
}

func (request *getPeriodsRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gpsOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		logger.SetHandler("get-periods")

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *getPeriodsRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetPeriodsHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gpsRequest == nil {
		gpsRequest = new(getPeriodsRequest)
	}

	err := gpsRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("get_periods_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer gpsRequest.finish()

	return gpsRequest.process(ctx, req)
}

func (request *getPeriodsRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	err := request.prepareRequest(req)
	if err != nil {
		return req.NewErrorResponse(err), nil
	}

	getPeriods := usecases.NewPeriodsGetter(request.periodRepo)

	userPeriods, nextKey, err := getPeriods(ctx, request.username, request.StartKey, request.PageSize, request.Active)
	if err != nil {
		request.err = err
		logger.Error("get_periods_failed", request.err, req)

		return req.NewErrorResponse(err), nil
	}

	res := &getPeriodsResponse{
		Periods: userPeriods,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, res), nil
}

func (request *getPeriodsRequest) prepareRequest(req *apigateway.Request) error {
	var err error

	request.username, err = apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_user_email_from_context_failed", err, req)

		return err
	}

	err = validate.Email(request.username)
	if err != nil {
		logger.Error("invalid_username", err, models.Any("user_data", map[string]interface{}{
			"s_username": request.username,
		}))

		return err
	}

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		logger.Error("get_request_params_failed", err, req)

		return err
	}

	val, _ := req.QueryStringParameters["active"]
	if strings.EqualFold(val, "true") {
		request.Active = true
	}

	return nil
}
