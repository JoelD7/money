package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var cpRequest *CreatePeriodRequest
var cpOnce sync.Once

type CreatePeriodRequest struct {
	startingTime   time.Time
	err            error
	PeriodRepo     period.Repository
	CacheManager   cache.IncomePeriodCacheManager
	SavingGoalRepo savingoal.Repository
	SavingsRepo    savings.Repository
}

func (request *CreatePeriodRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	cpOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.PeriodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}

		request.SavingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.SavingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.CacheManager = cache.NewRedisCache()
		logger.SetHandler("create-period")
	})
	request.startingTime = time.Now()

	return err
}

func (request *CreatePeriodRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreatePeriodHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if cpRequest == nil {
		cpRequest = new(CreatePeriodRequest)
	}

	err := cpRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("create_period_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer cpRequest.finish()

	return cpRequest.Process(ctx, req)
}

func (request *CreatePeriodRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	_, err := req.GetIdempotenceyKeyFromHeader()
	if err != nil {
		request.err = err
		logger.Error("http_request_validation_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	periodModel, err := request.validateCreateRequestBody(req)
	if err != nil {
		request.err = err
		logger.Error("validate_request_body_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	createPeriod := usecases.NewPeriodCreator(request.PeriodRepo, request.CacheManager, request.SavingGoalRepo, request.SavingsRepo)

	createdPeriod, err := createPeriod(ctx, username, periodModel)
	if err != nil {
		request.err = err
		logger.Error("create_period_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, createdPeriod), nil
}

func (request *CreatePeriodRequest) validateCreateRequestBody(req *apigateway.Request) (*models.Period, error) {
	p := new(models.Period)

	err := json.Unmarshal([]byte(req.Body), p)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, models.ErrInvalidRequestBody)
	}

	if p.Name == nil || p.Name != nil && *p.Name == "" {
		return nil, models.ErrMissingPeriodName
	}

	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	return p, nil
}
