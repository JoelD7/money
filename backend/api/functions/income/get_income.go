package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

type incomeGetRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	incomeRepo   income.Repository
}

var once sync.Once
var request *incomeGetRequest

func (request *incomeGetRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	once.Do(func() {
		request.log = log
		dynamoClient := dynamo.InitClient(ctx)

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig.IncomeTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *incomeGetRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getIncomeHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if request == nil {
		request = new(incomeGetRequest)
	}

	request.init(ctx, log, envConfig)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *incomeGetRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	incomeID, ok := req.PathParameters["incomeID"]
	if !ok || incomeID == "" {
		request.err = models.ErrMissingIncomeID

		request.log.Error("missing_income_id", nil, []models.LoggerObject{req})
		return req.NewErrorResponse(models.ErrMissingIncomeID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err

		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}

	getIncome := usecases.NewIncomeGetter(request.incomeRepo)

	userIncome, err := getIncome(ctx, username, incomeID)
	if err != nil {
		request.err = err

		request.log.Error("get_income_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, userIncome), nil
}
