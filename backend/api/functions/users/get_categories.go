package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type getCategoriesRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *getCategoriesRequest) init() {
	dynamoClient := initDynamoClient()

	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLoggerWithHandler("get-categories")
}

func (request *getCategoriesRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getCategoriesHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getCategoriesRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getCategoriesRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := getUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	getCategories := usecases.NewCategoriesGetter(request.userRepo)

	categories, err := getCategories(ctx, username)
	if err != nil {
		request.err = err
		request.log.Error("get_categories_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, categories), nil
}
