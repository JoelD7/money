package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type requestLogoutHandler struct {
	RefreshToken string `json:"refresh_token,omitempty"`

	log                 logger.LogAPI
	startingTime        time.Time
	err                 error
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
}

func logoutHandler(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestLogoutHandler{}

	req.initLogoutHandler()
	defer req.finish()

	return req.processLogout(ctx, request)
}

func (req *requestLogoutHandler) initLogoutHandler() {
	dynamoClient := initDynamoClient()

	req.userRepo = users.NewDynamoRepository(dynamoClient)
	req.invalidTokenManager = cache.NewRedisCache()
	req.startingTime = time.Now()
	req.log = logger.NewLoggerWithHandler("logout")
}

func (req *requestLogoutHandler) finish() {
	defer func() {
		err := req.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLogoutHandler) processLogout(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	var err error

	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		req.log.Error("getting_refresh_token_cookie_failed", err, nil)

		return getErrorResponse(err)
	}

	logout := usecases.NewUserLogout(req.userRepo, req.invalidTokenManager, req.log)
	err = logout(ctx, req.RefreshToken)
	if err != nil {
		req.err = err
		req.log.Error("logout_failed", err, nil)

		return getErrorResponse(err)
	}

	setCookieHeader := map[string]string{
		"Set-Cookie": getExpiredRefreshTokenCookie(),
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Headers:    setCookieHeader,
	}, nil
}

func getExpiredRefreshTokenCookie() string {
	t := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)

	return fmt.Sprintf("%s=; Path=/; Expires=%s; Secure; HttpOnly", refreshTokenCookieName, t.Format(time.RFC1123))
}
