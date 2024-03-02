package main

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errNoCategoryIDInPath = errors.New("no category id in path")
)

type updateCategoryRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *updateCategoryRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("update-category")
}

func (request *updateCategoryRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updateCategoryHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updateCategoryRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updateCategoryRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	categoryID, ok := req.PathParameters["categoryID"]
	if !ok {
		request.err = errNoCategoryIDInPath
		request.log.Error("get_category_id_from_path_failed", errNoCategoryIDInPath, []models.LoggerObject{req})

		return req.NewErrorResponse(errNoCategoryIDInPath), nil
	}

	requestCategory, err := validateRequestBody(req)
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

	updateCategory := usecases.NewCategoryUpdater(request.userRepo)

	err = updateCategory(ctx, username, categoryID, requestCategory)
	if err != nil {
		request.err = err
		request.log.Error("update_category_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func validateRequestBody(req *apigateway.Request) (*models.Category, error) {
	requestCategory := new(models.Category)

	err := json.Unmarshal([]byte(req.Body), requestCategory)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	if requestCategory.Name != nil && *requestCategory.Name == "" {
		return nil, models.ErrMissingCategoryName
	}

	if requestCategory.Budget != nil && (*requestCategory.Budget < 0 || *requestCategory.Budget >= math.MaxFloat64) {
		return nil, models.ErrInvalidBudget
	}

	return requestCategory, nil
}
