package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
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
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *createCategoryRequest) init(ctx context.Context, log logger.LogAPI) {
	ccOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo = users.NewDynamoRepository(dynamoClient)
		request.log = log
		request.log.SetHandler("create-category")
	})
	request.startingTime = time.Now()
}

func (request *createCategoryRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateCategoryHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if ccRequest == nil {
		ccRequest = new(createCategoryRequest)
	}

	ccRequest.init(ctx, log)
	defer ccRequest.finish()

	return ccRequest.process(ctx, req)
}

func (request *createCategoryRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	category, err := validateCreateCategoryRequestBody(req)
	if err != nil {
		request.log.Error("request_body_validation_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createCategory := usecases.NewCategoryCreator(request.userRepo)

	err = createCategory(ctx, username, category)
	if err != nil {
		request.err = err
		request.log.Error("create_category_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
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
