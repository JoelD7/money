package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gcRequest *getCategoriesRequest
var gcOnce sync.Once

type getCategoriesRequest struct {
	startingTime time.Time
	err          error
	userRepo     users.Repository
}

func (request *getCategoriesRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gcOnce.Do(func() {
		logger.SetHandler("get-categories")
		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *getCategoriesRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetCategoriesHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gcRequest == nil {
		gcRequest = new(getCategoriesRequest)
	}

	err := gcRequest.init(ctx, envConfig)
	if err != nil {
		gcRequest.err = err

		logger.Error("get_categories_init_failed", err, req)

		return req.NewErrorResponse(err), nil

	}
	defer gcRequest.finish()

	return gcRequest.process(ctx, req)
}

func (request *getCategoriesRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
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

	getCategories := usecases.NewCategoriesGetter(request.userRepo)

	categories, err := getCategories(ctx, username)
	if err != nil {
		request.err = err
		logger.Error("get_categories_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, categories), nil
}
