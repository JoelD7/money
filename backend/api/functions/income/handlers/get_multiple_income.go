package handlers

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"strconv"
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
}

type MultipleIncomeResponse struct {
	Income  []*models.Income `json:"income"`
	Periods []string         `json:"periods"`
	NextKey string           `json:"next_key"`
}

func (request *GetMultipleIncomeRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gmiOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.IncomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig.IncomeTable, envConfig.PeriodUserIncomeIndex)
		if err != nil {
			return
		}
		request.Log = log

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
		request.log.Error("init_failed", err, []models.LoggerObject{req})

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
		request.Log.Error("get_username_failed", nil, []models.LoggerObject{req})

		return err
	}

	request.Username = username

	request.StartKey, request.PageSize, err = getRequestQueryParams(req)
	if err != nil {
		request.Log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return err
	}

	return nil
}

func getRequestQueryParams(req *apigateway.Request) (string, int, error) {
	pageSizeParam := 0
	var err error

	if req.QueryStringParameters["page_size"] != "" {
		pageSizeParam, err = strconv.Atoi(req.QueryStringParameters["page_size"])
	}

	if err != nil || pageSizeParam < 0 {
		return "", 0, fmt.Errorf("%w: %v", models.ErrInvalidPageSize, err)
	}

	return req.QueryStringParameters["start_key"], pageSizeParam, nil
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

		request.Log.Error("missing_period", nil, []models.LoggerObject{req})
		return req.NewErrorResponse(models.ErrMissingPeriod), nil
	}

	getIncomeByPeriod := usecases.NewIncomeByPeriodGetter(request.IncomeRepo, request.CacheManager)

	userIncome, nextKey, incomePeriods, err := getIncomeByPeriod(ctx, request.Username, period, request.StartKey, request.PageSize)
	if err != nil {
		request.err = err
		request.Log.Error("get_income_by_period_failed", err, []models.LoggerObject{req})
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

	userIncome, nextKey, incomePeriods, err := getAllIncome(ctx, request.Username, request.StartKey, request.PageSize)
	if err != nil {
		request.err = err
		request.Log.Error("get_all_income_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}

	response := &MultipleIncomeResponse{
		Income:  userIncome,
		Periods: incomePeriods,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}
