package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var gmiRequest *getMultipleIncomeRequest
var gmiOnce sync.Once

type getMultipleIncomeRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	username     string
	startKey     string
	pageSize     int
	incomeRepo   income.Repository
}

type multipleIncomeResponse struct {
	Income  []*models.Income `json:"income"`
	NextKey string           `json:"next_key"`
}

func (request *getMultipleIncomeRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gmiOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig.IncomeTable)
		if err != nil {
			return
		}
		request.log = log
	})
	request.startingTime = time.Now()

	return err
}

func (request *getMultipleIncomeRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getMultipleIncomeHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gmiRequest == nil {
		gmiRequest = new(getMultipleIncomeRequest)
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

	return gmiRequest.routeToHandlers(ctx, req)
}

func (request *getMultipleIncomeRequest) prepareRequest(req *apigateway.Request) error {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_failed", nil, []models.LoggerObject{req})

		return err
	}

	request.username = username

	request.startKey, request.pageSize, err = getRequestQueryParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

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

func (request *getMultipleIncomeRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period, ok := req.QueryStringParameters["period"]
	if ok {
		return request.getIncomeByPeriod(ctx, req, period)
	}

	return request.getAllIncome(ctx, req)
}

func (request *getMultipleIncomeRequest) getIncomeByPeriod(ctx context.Context, req *apigateway.Request, period string) (*apigateway.Response, error) {
	if period == "" {
		request.err = models.ErrMissingPeriod

		request.log.Error("missing_period", nil, []models.LoggerObject{req})
		return req.NewErrorResponse(models.ErrMissingPeriod), nil
	}

	getIncomeByPeriod := usecases.NewIncomeByPeriodGetter(request.incomeRepo)

	userIncome, nextKey, err := getIncomeByPeriod(ctx, request.username, period, request.startKey, request.pageSize)
	if err != nil {
		request.err = err
		request.log.Error("get_income_by_period_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}

	response := &multipleIncomeResponse{
		Income:  userIncome,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}

func (request *getMultipleIncomeRequest) getAllIncome(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getAllIncome := usecases.NewAllIncomeGetter(request.incomeRepo)

	userIncome, nextKey, err := getAllIncome(ctx, request.username, request.startKey, request.pageSize)
	if err != nil {
		request.err = err
		request.log.Error("get_all_income_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}

	response := &multipleIncomeResponse{
		Income:  userIncome,
		NextKey: nextKey,
	}

	return req.NewJSONResponse(http.StatusOK, response), nil
}
