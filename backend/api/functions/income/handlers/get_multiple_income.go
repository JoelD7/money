package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gmiRequest *GetMultipleIncomeRequest
var gmiOnce sync.Once

type GetMultipleIncomeRequest struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	Username     string
	StartKey     string
	PageSize     int
	IncomeRepo   income.Repository
	CacheManager cache.IncomePeriodCacheManager
	*models.QueryParameters
}

type MultipleIncomeResponse struct {
	Income  []*models.Income `json:"income"`
	Periods []string         `json:"periods"`
	NextKey string           `json:"next_key"`
}

func (request *GetMultipleIncomeRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gmiOnce.Do(func() {
		request.Log = log

		dynamoClient := dynamo.InitClient(ctx)

		request.IncomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.CacheManager = cache.NewRedisCache()
	})
	request.startingTime = time.Now()

	return err
}

func (request *GetMultipleIncomeRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetMultipleIncomeHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gmiRequest == nil {
		gmiRequest = new(GetMultipleIncomeRequest)
	}

	err := gmiRequest.init(ctx, log, envConfig)
	if err != nil {
		gmiRequest.Log.Error("init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	defer gmiRequest.finish()

	err = gmiRequest.prepareRequest(req)
	if err != nil {
		return req.NewErrorResponse(err), nil
	}

	return gmiRequest.RouteToHandlers(ctx, req)
}

func (request *GetMultipleIncomeRequest) prepareRequest(req *apigateway.Request) error {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.Log.Error("get_username_failed", nil, req)

		return err
	}

	request.Username = username

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		request.Log.Error("get_request_params_failed", err, req)

		return err
	}

	return nil
}

func (request *GetMultipleIncomeRequest) RouteToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period, ok := req.QueryStringParameters["period"]
	if ok {
		return request.GetIncomeByPeriod(ctx, req, period)
	}

	return request.getAllIncome(ctx, req)
}

func (request *GetMultipleIncomeRequest) GetIncomeByPeriod(ctx context.Context, req *apigateway.Request, period string) (*apigateway.Response, error) {
	if period == "" {
		request.err = models.ErrMissingPeriod

		request.Log.Error("missing_period", nil, req)
		return req.NewErrorResponse(models.ErrMissingPeriod), nil
	}

	getIncomeByPeriod := usecases.NewIncomeByPeriodGetter(request.IncomeRepo, request.CacheManager)

	userIncome, nextKey, incomePeriods, err := getIncomeByPeriod(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.err = err
		request.Log.Error("get_income_by_period_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	response := &MultipleIncomeResponse{
		Income:  userIncome,
		Periods: incomePeriods,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}

func (request *GetMultipleIncomeRequest) getAllIncome(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getAllIncome := usecases.NewAllIncomeGetter(request.IncomeRepo, request.CacheManager)

	userIncome, nextKey, incomePeriods, err := getAllIncome(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.err = err
		request.Log.Error("get_all_income_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	response := &MultipleIncomeResponse{
		Income:  userIncome,
		Periods: incomePeriods,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}
