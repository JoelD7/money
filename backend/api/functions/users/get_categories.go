package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gcRequest *getCategoriesRequest
var gcOnce sync.Once

type getCategoriesRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *getCategoriesRequest) init(log logger.LogAPI) {
	gcOnce.Do(func() {
		dynamoClient := initDynamoClient()

		request.userRepo = users.NewDynamoRepository(dynamoClient)
		request.log = log
		request.log.SetHandler("get-categories")
	})
	request.startingTime = time.Now()
}

func (request *getCategoriesRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getCategoriesHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if gcRequest == nil {
		gcRequest = new(getCategoriesRequest)
	}

	gcRequest.init(log)
	defer gcRequest.finish()

	return gcRequest.process(ctx, req)
}

func (request *getCategoriesRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
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

	getCategories := usecases.NewCategoriesGetter(request.userRepo)

	categories, err := getCategories(ctx, username)
	if err != nil {
		request.err = err
		request.log.Error("get_categories_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, categories), nil
}
