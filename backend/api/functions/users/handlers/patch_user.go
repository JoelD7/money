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
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	puRequest *patchUserRequest
	puOnce    sync.Once
)

type patchUserRequest struct {
	startingTime time.Time
	err          error
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *patchUserRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	puOnce.Do(func() {
		logger.SetHandler("patch-user")
		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *patchUserRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func PatchUserHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if puRequest == nil {
		puRequest = new(patchUserRequest)
	}

	err := puRequest.init(ctx, envConfig)
	if err != nil {
		puRequest.err = err

		logger.Error("update_category_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer puRequest.finish()

	return puRequest.process(ctx, req)
}

func (request *patchUserRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, ok := req.PathParameters["username"]
	if !ok {
		err := fmt.Errorf("missing username")
		logger.Error("missing_username", err, req)

		return req.NewErrorResponse(err), nil
	}

	err := validate.Email(username)
	if err != nil {
		logger.Error("invalid_username", err, req)

		return req.NewErrorResponse(err), nil
	}

	user, err := validateUserBody(req)
	if err != nil {
		logger.Error("invalid_user_body", err, req)

		return req.NewErrorResponse(err), nil
	}

	user.Username = username

	patchUser := usecases.NewUserPatcher(request.userRepo, request.periodRepo)

	err = patchUser(ctx, username, user)
	if err != nil {
		request.err = err
		logger.Error("patch_user_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, nil), nil
}

func validateUserBody(req *apigateway.Request) (*models.User, error) {
	var user models.User

	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidRequestBody)
	}

	return &user, nil
}
