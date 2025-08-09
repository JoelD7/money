package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
)

var gmiRequest *GetMultipleIncomeRequest
var gmiOnce sync.Once

type GetMultipleIncomeRequest struct {
	startingTime time.Time
	err          error
	Username     string
	StartKey     string
	PageSize     int
	IncomeRepo   income.Repository
	CacheManager cache.IncomePeriodCacheManager
	PeriodRepo   period.Repository
	*models.QueryParameters
}

type MultipleIncomeResponse struct {
	Income  []*models.Income `json:"income"`
	Periods []string         `json:"periods"`
	NextKey string           `json:"next_key"`
}

func (request *GetMultipleIncomeRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gmiOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.IncomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.PeriodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}

		request.CacheManager = cache.NewRedisCache()
	})
	request.startingTime = time.Now()

	return err
}

func (request *GetMultipleIncomeRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetMultipleIncomeHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gmiRequest == nil {
		gmiRequest = new(GetMultipleIncomeRequest)
	}

	err := gmiRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_failed", err, req)

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
		logger.Error("get_username_failed", nil, req)

		return err
	}

	request.Username = username

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		logger.Error("get_request_params_failed", err, req)

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

		logger.Error("missing_period", nil, req)
		return req.NewErrorResponse(models.ErrMissingPeriod), nil
	}

	getIncomeByPeriod := usecases.NewIncomeByPeriodGetter(request.IncomeRepo, request.CacheManager, request.PeriodRepo)

	userIncome, nextKey, incomePeriods, err := getIncomeByPeriod(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.err = err
		logger.Error("get_income_by_period_failed", err, req)
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
	getAllIncome := usecases.NewAllIncomeGetter(request.IncomeRepo, request.CacheManager, request.PeriodRepo)

	userIncome, nextKey, incomePeriods, err := getAllIncome(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.err = err
		logger.Error("get_all_income_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	response := &MultipleIncomeResponse{
		Income:  userIncome,
		Periods: incomePeriods,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}
