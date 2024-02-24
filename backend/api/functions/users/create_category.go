package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"math"
	"net/http"
	"time"
)

type createCategoryRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *createCategoryRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("create-category")
}

func (request *createCategoryRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createCategoryHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(createCategoryRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *createCategoryRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	category, err := validateCreateCategoryRequestBody(req)
	if err != nil {
		request.log.Error("request_body_validation_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	createCategory := usecases.NewCategoryCreator(request.userRepo)

	err = createCategory(ctx, username, category)
	if err != nil {
		request.err = err
		request.log.Error("create_category_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
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
