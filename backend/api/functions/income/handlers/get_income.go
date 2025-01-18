package handlers

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
	startingTime time.Time
	err          error
	incomeRepo   income.Repository
}

var once sync.Once
var request *incomeGetRequest

func (request *incomeGetRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	once.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *incomeGetRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetIncomeHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if request == nil {
		request = new(incomeGetRequest)
	}

	err := request.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	defer request.finish()

	return request.process(ctx, req)
}

func (request *incomeGetRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	incomeID, ok := req.PathParameters["incomeID"]
	if !ok || incomeID == "" {
		request.err = models.ErrMissingIncomeID

		logger.Error("missing_income_id", nil, req)
		return req.NewErrorResponse(models.ErrMissingIncomeID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err

		logger.Error("get_username_from_context_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	getIncome := usecases.NewIncomeGetter(request.incomeRepo)

	userIncome, err := getIncome(ctx, username, incomeID)
	if err != nil {
		request.err = err

		logger.Error("get_income_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, userIncome), nil
}
