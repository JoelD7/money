package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type getMultipleIncomeRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	username     string
	incomeRepo   income.Repository
}

type multipleIncomeResponse struct {
	Income  []*models.Income `json:"income"`
	NextKey string           `json:"next_key"`
}

func (request *getMultipleIncomeRequest) init() {
	dynamoClient := initDynamoClient()

	request.incomeRepo = income.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getMultipleIncomeRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getMultipleIncomeHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getMultipleIncomeRequest)

	request.init()
	defer request.finish()

	err := request.prepareRequest(req)
	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	return request.routeToHandlers(ctx, req)
}

func (request *getMultipleIncomeRequest) prepareRequest(req *apigateway.Request) error {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_failed", nil, []models.LoggerObject{req})

		return err
	}

	request.username = username

	return nil
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
		return apigateway.NewErrorResponse(models.ErrMissingPeriod), nil
	}

	getIncomeByPeriod := usecases.NewIncomeByPeriodGetter(request.incomeRepo)

	userIncome, nextKey, err := getIncomeByPeriod(ctx, request.username, period, request.startKey, request.pageSize)
	if err != nil {
		request.err = err
		request.log.Error("get_income_by_period_failed", err, []models.LoggerObject{req})
		return apigateway.NewErrorResponse(err), nil
	}

	response := &multipleIncomeResponse{
		Income:  userIncome,
		NextKey: nextKey,
	}

	return apigateway.NewJSONResponse(http.StatusOK, response), nil
}

func (request *getMultipleIncomeRequest) getAllIncome(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {

}
