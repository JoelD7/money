package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	ciRequest *createIncomeRequest
	ciOnce    sync.Once
)

type createIncomeRequest struct {
	startingTime     time.Time
	err              error
	incomeRepo       income.Repository
	periodRepo       period.Repository
	idempotenceCache cache.IdempotenceCacheManager
}

func (request *createIncomeRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	ciOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
		request.idempotenceCache = cache.NewRedisCache()
	})
	request.startingTime = time.Now()

	return err
}

func (request *createIncomeRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateIncomeHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if ciRequest == nil {
		ciRequest = new(createIncomeRequest)
	}

	err := ciRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	defer ciRequest.finish()

	return ciRequest.process(ctx, req)
}

func (request *createIncomeRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	idempotencyKey, err := req.GetIdempotenceyKeyFromHeader()
	if err != nil {
		request.err = err
		logger.Error("http_request_validation_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	reqIncome, err := validateCreateIncomeBody(req)
	if err != nil {
		request.err = err
		logger.Error("validate_create_income_body_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	createIncome := usecases.NewIncomeCreator(request.incomeRepo, request.periodRepo, request.idempotenceCache)

	newIncome, err := createIncome(ctx, username, idempotencyKey, reqIncome)
	if err != nil {
		request.err = err
		logger.Error("create_income_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, newIncome), nil
}

func validateCreateIncomeBody(req *apigateway.Request) (*models.Income, error) {
	reqIncome := new(models.Income)

	err := json.Unmarshal([]byte(req.Body), &reqIncome)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	if reqIncome.Amount == nil {
		return nil, models.ErrMissingAmount
	}

	if reqIncome.Name == nil {
		return nil, models.ErrMissingName
	}

	if reqIncome.PeriodID == nil {
		return nil, models.ErrMissingPeriod
	}

	err = validate.Amount(reqIncome.Amount)
	if err != nil {
		return nil, err
	}

	return reqIncome, nil
}
