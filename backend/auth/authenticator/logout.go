package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var logoutRequest *requestLogoutHandler
var logoutOnce sync.Once

type requestLogoutHandler struct {
	startingTime        time.Time
	err                 error
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
}

func logoutHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if logoutRequest == nil {
		logoutRequest = new(requestLogoutHandler)
	}

	err := logoutRequest.initLogoutHandler(ctx, envConfig)
	if err != nil {
		logoutRequest.err = err

		logger.Error("logout_init_failed", err, request)

		return request.NewErrorResponse(err), nil
	}
	defer logoutRequest.finish()

	return logoutRequest.processLogout(ctx, request)
}

func (req *requestLogoutHandler) initLogoutHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	logoutOnce.Do(func() {
		logger.SetHandler("logout")
		dynamoClient := dynamo.InitClient(ctx)

		req.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
		req.invalidTokenManager = cache.NewRedisCache()
	})
	req.startingTime = time.Now()
	req.err = nil

	return err
}

func (req *requestLogoutHandler) finish() {
	logger.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLogoutHandler) processLogout(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	var err error

	credentials, err := validateRequestBody(request)
	if err != nil {
		req.err = err
		logger.Error("logout_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	logout := usecases.NewUserLogout(req.userRepo, req.invalidTokenManager)

	err = logout(ctx, credentials.Username)
	if errors.Is(err, models.ErrUserNotFound) {
		req.err = err
		logger.Error("logout_failed", err, nil)

		return request.NewErrorResponse(errUserNotFound), nil
	}

	if err != nil {
		req.err = err
		logger.Error("logout_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	return request.NewJSONResponse(http.StatusOK, nil, apigateway.Header{
		Key:   "Set-Cookie",
		Value: getExpiredRefreshTokenCookie(),
	}), nil
}

func validateRequestBody(request *apigateway.Request) (*Credentials, error) {
	var credentials *Credentials

	err := json.Unmarshal([]byte(request.Body), &credentials)
	if err != nil {
		err = fmt.Errorf("%w: %v", models.ErrInvalidRequestBody, err)
		req.err = err
		logger.Error("unmarshal_credentials_failed", err, nil)

		return nil, err
	}

	if credentials.Username == "" {
		return nil, models.ErrMissingUsername
	}

	return credentials, nil
}

func getExpiredRefreshTokenCookie() string {
	t := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)

	return fmt.Sprintf("%s=; Path=/; Expires=%s; Secure; HttpOnly; SameSite=None", refreshTokenCookieName, t.Format(time.RFC1123))
}
