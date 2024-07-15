package handlers

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	errNoCategoryIDInPath = errors.New("no category id in path")
	ucRequest             *updateCategoryRequest
	ucOnce                sync.Once
)

type updateCategoryRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *updateCategoryRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	ucOnce.Do(func() {
		request.log = log
		request.log.SetHandler("update-category")
		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, usersTableName)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *updateCategoryRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func UpdateCategoryHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if ucRequest == nil {
		ucRequest = new(updateCategoryRequest)
	}

	err := ucRequest.init(ctx, log)
	if err != nil {
		ucRequest.err = err

		log.Error("update_category_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}
	defer ucRequest.finish()

	return ucRequest.process(ctx, req)
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
