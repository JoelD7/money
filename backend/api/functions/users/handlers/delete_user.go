package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
)

var (
	duRequest *deleteUserRequest
	duOnce    sync.Once
)

type deleteUserRequest struct {
	userRepo users.Repository

	startingTime time.Time
	err          error
}

func (dur *deleteUserRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	duOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		dur.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
	})

	dur.startingTime = time.Now()

	return err
}

func (dur *deleteUserRequest) finish() {
	logger.LogLambdaTime(dur.startingTime, dur.err, recover())
}

func DeleteUserHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if duRequest == nil {
		duRequest = new(deleteUserRequest)
	}

	err := duRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("delete_user_init_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	defer duRequest.finish()

	return duRequest.process(ctx, req)
}

func (dur *deleteUserRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	pathUsername := req.PathParameters["username"]

	authorizerUsername, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	if pathUsername != authorizerUsername {
		return req.NewErrorResponse(models.ErrUsernameDeleteMismatch), nil
	}

	deleteUser := usecases.NewUserDeleter(dur.userRepo)

	err = deleteUser(ctx, pathUsername)
	if err != nil {
		logger.Error("delete_user_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusNoContent, nil), nil
}
