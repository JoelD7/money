package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var logoutRequest *requestLogoutHandler
var logoutOnce sync.Once

type requestLogoutHandler struct {
	RefreshToken string `json:"refresh_token,omitempty"`

	log                 logger.LogAPI
	startingTime        time.Time
	err                 error
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
}

func logoutHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	if logoutRequest == nil {
		logoutRequest = new(requestLogoutHandler)
	}

	logoutRequest.initLogoutHandler(log)
	defer logoutRequest.finish()

	return logoutRequest.processLogout(ctx, request)
}

func (req *requestLogoutHandler) initLogoutHandler(log logger.LogAPI) {
	logoutOnce.Do(func() {
		dynamoClient := initDynamoClient()

		req.userRepo = users.NewDynamoRepository(dynamoClient)
		req.invalidTokenManager = cache.NewRedisCache()
		req.log = log
		req.log.SetHandler("logout")
	})
	req.startingTime = time.Now()
}

func (req *requestLogoutHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLogoutHandler) processLogout(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	var err error

	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		req.log.Error("getting_refresh_token_cookie_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	logout := usecases.NewUserLogout(req.userRepo, req.invalidTokenManager, req.log)

	err = logout(ctx, req.RefreshToken)
	if errors.Is(err, models.ErrUserNotFound) {
		req.err = err
		req.log.Error("logout_failed", err, nil)

		return request.NewErrorResponse(errUserNotFound), nil
	}

	if err != nil {
		req.err = err
		req.log.Error("logout_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	return request.NewJSONResponse(http.StatusOK, nil, apigateway.Header{
		Key:   "Set-Cookie",
		Value: getExpiredRefreshTokenCookie(),
	}), nil
}

func getExpiredRefreshTokenCookie() string {
	t := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)

	return fmt.Sprintf("%s=; Path=/; Expires=%s; Secure; HttpOnly; SameSite=None", refreshTokenCookieName, t.Format(time.RFC1123))
}
