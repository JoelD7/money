package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	csRequest *createSavingRequest
	csOnce    sync.Once
)

type createSavingRequest struct {
	startingTime     time.Time
	err              error
	savingsRepo      savings.Repository
	userRepo         users.Repository
	periodRepo       period.Repository
	idempotenceCache cache.IdempotenceCacheManager
}

func (request *createSavingRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	csOnce.Do(func() {
		logger.SetHandler("create-saving")
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}

		request.idempotenceCache = cache.NewRedisCache()
		request.idempotenceCache.SetTTL(envConfig.IdempotencyKeyCacheTTLSeconds)
	})
	request.startingTime = time.Now()

	return err
}

func (request *createSavingRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateSavingHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if csRequest == nil {
		csRequest = new(createSavingRequest)
	}

	err := csRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_create_saving_request_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer csRequest.finish()

	return csRequest.process(ctx, req)
}

func (request *createSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	idempotencyKey, err := req.GetIdempotenceyKeyFromHeader()
	if err != nil {
		request.err = err
		logger.Error("http_request_validation_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	userSaving, err := validateBody(req)
	if err != nil {
		logger.Error("validate_request_body_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	createSaving := usecases.NewSavingCreator(request.savingsRepo, request.periodRepo, request.idempotenceCache)

	saving, err := createSaving(ctx, username, idempotencyKey, userSaving)
	if err != nil {
		logger.Error("create_saving_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, saving), nil
}

func validateBody(req *apigateway.Request) (*models.Saving, error) {
	userSaving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), userSaving)
	if err != nil {
		return nil, models.ErrInvalidRequestBody
	}

	if userSaving.Amount == nil || (userSaving.Amount != nil && *userSaving.Amount == 0) {
		return nil, models.ErrMissingAmount
	}

	err = validate.Amount(userSaving.Amount)
	if err != nil {
		return nil, models.ErrInvalidSavingAmount
	}

	if userSaving.Period == nil || (userSaving.Period != nil && *userSaving.Period == "") {
		return nil, models.ErrMissingPeriod
	}

	return userSaving, nil
}
