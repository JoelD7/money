package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
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

	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     *users.Repository
	cacheRepo    *cache.Repository
}

func logoutHandler(request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestLogoutHandler{
		log: logger.NewLoggerWithHandler("logout"),
	}

	req.initLogoutHandler()
	defer req.finish()

	return req.processLogout(request)
}

func (req *requestLogoutHandler) initLogoutHandler() {
	dynamoClient := initDynamoClient()

	dynamoUserRepository := users.NewDynamoRepository(dynamoClient)

	req.userRepo = users.NewRepository(dynamoUserRepository)
	redisRepository := cache.NewRepository(cache.NewRedisCache())
	req.cacheRepo = redisRepository
	req.startingTime = time.Now()
	req.log = logger.NewLogger()
}

func (req *requestLogoutHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLogoutHandler) processLogout(request *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	var err error

	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		req.log.Error("getting_refresh_token_cookie_failed", err, nil)

		return getErrorResponse(err)
	}

	logout := usecases.NewUserLogout(req.userRepo, req.cacheRepo, req.log)
	err = logout(ctx, req.RefreshToken)
	if errors.Is(err, models.ErrInvalidToken) {
		req.err = err
		req.log.Error("token_payload_parse_failed", err, nil)

		return getErrorResponse(err)
	}

	if err != nil {
		req.err = err
		req.log.Error("token_payload_parse_failed", err, nil)

		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
}
