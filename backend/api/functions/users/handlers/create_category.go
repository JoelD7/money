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
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"math"
	"net/http"
	"sync"
	"time"
)

var ccRequest *createCategoryRequest
var ccOnce sync.Once

type createCategoryRequest struct {
	startingTime     time.Time
	err              error
	userRepo         users.Repository
	idempotenceCache cache.IdempotenceCacheManager
}

func (request *createCategoryRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	ccOnce.Do(func() {
		logger.SetHandler("create-category")
		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}

		request.idempotenceCache = cache.NewRedisCache()
		request.idempotenceCache.SetTTL(envConfig.IdempotencyKeyCacheTTLSeconds)
	})
	request.startingTime = time.Now()

	return err
}

func (request *createCategoryRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateCategoryHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if ccRequest == nil {
		ccRequest = new(createCategoryRequest)
	}

	err := ccRequest.init(ctx, envConfig)
	if err != nil {
		ccRequest.err = err

		logger.Error("create_category_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer ccRequest.finish()

	return ccRequest.process(ctx, req)
}

func (request *createCategoryRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	idempotencyKey, err := req.GetIdempotenceyKeyFromHeader()
	if err != nil {
		request.err = err
		logger.Error("http_request_validation_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	category, err := validateCreateCategoryRequestBody(req)
	if err != nil {
		logger.Error("request_body_validation_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		logger.Error("get_user_email_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		logger.Error("invalid_username", err, req)

		return req.NewErrorResponse(err), nil
	}

	createCategory := usecases.NewCategoryCreator(request.userRepo, request.idempotenceCache)

	err = createCategory(ctx, username, idempotencyKey, category)
	if err != nil {
		request.err = err
		logger.Error("create_category_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, nil), nil
}

func validateCreateCategoryRequestBody(req *apigateway.Request) (*models.Category, error) {
	category := new(models.Category)

	err := json.Unmarshal([]byte(req.Body), category)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	if category.Name == nil || *category.Name == "" {
		return nil, models.ErrMissingCategoryName
	}

	if category.Budget == nil {
		return nil, models.ErrMissingCategoryBudget
	}

	if *category.Budget < 0 || *category.Budget >= math.MaxFloat64 {
		return nil, models.ErrInvalidBudget
	}

	if category.Color == nil || *category.Color == "" {
		return nil, models.ErrMissingCategoryColor
	}

	return category, nil
}
