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

type getIncomeRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	incomeRepo   income.Repository
}

func (request *getIncomeRequest) init() {
	dynamoClient := initDynamoClient()

	request.incomeRepo = income.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getIncomeRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getIncomeHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getIncomeRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getIncomeRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	incomeID, ok := req.PathParameters["incomeID"]
	if !ok || incomeID == "" {
		request.err = models.ErrMissingIncomeID

		request.log.Error("missing_income_id", nil, []models.LoggerObject{req})
		return apigateway.NewErrorResponse(models.ErrMissingIncomeID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err

		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})
		return apigateway.NewErrorResponse(err), nil
	}

	getIncome := usecases.NewIncomeGetter(request.incomeRepo)

	userIncome, err := getIncome(ctx, username, incomeID)
	if err != nil {
		request.err = err

		request.log.Error("get_income_failed", err, []models.LoggerObject{req})
		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, userIncome), nil
}
