package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	once sync.Once
	req  *requestTokenHandler
)

type requestTokenHandler struct {
	RefreshToken string `json:"refresh_token,omitempty"`

	log                 logger.LogAPI
	startingTime        time.Time
	err                 error
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
	secretsManager      secrets.SecretManager
}

func tokenHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	if req == nil {
		req = new(requestTokenHandler)
	}

	req.initTokenHandler(log)
	defer req.finish()

	return req.processToken(ctx, request)
}

func (req *requestTokenHandler) initTokenHandler(log logger.LogAPI) {
	once.Do(func() {
		dynamoClient := initDynamoClient()

		req.userRepo = users.NewDynamoRepository(dynamoClient)
		req.invalidTokenManager = cache.NewRedisCache()
		req.secretsManager = secrets.NewAWSSecretManager()
		req.log = log
		req.log.SetHandler("token")
	})

	req.startingTime = time.Now()
}

func (req *requestTokenHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestTokenHandler) processToken(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	var err error
	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		req.log.Error("getting_refresh_token_cookie_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	validateRefreshToken := usecases.NewRefreshTokenValidator(req.userRepo, req.log)

	user, err := validateRefreshToken(ctx, req.RefreshToken)

	if err != nil && errors.Is(err, models.ErrInvalidToken) {
		req.err = err

		return req.handleValidationError(ctx, user, request)
	}

	if errors.Is(err, models.ErrUserNotFound) {
		return request.NewErrorResponse(errUserNotFound), nil
	}

	if err != nil {
		req.err = err

		return request.NewErrorResponse(err), nil
	}

	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager, req.log)

	accessToken, refreshToken, err := generateTokens(ctx, user)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	response := &accessTokenResponse{accessToken.Value}

	data, err := json.Marshal(response)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	cookieStr := getRefreshTokenCookieStr(refreshToken.Value, refreshToken.Expiration)

	req.log.Info("new_tokens_issued_successfully", []models.LoggerObject{user})

	return request.NewJSONResponse(http.StatusOK, string(data), apigateway.Header{Key: "Set-Cookie", Value: cookieStr}), nil
}

func (req *requestTokenHandler) handleValidationError(ctx context.Context, user *models.User, request *apigateway.Request) (*apigateway.Response, error) {
	invalidateTokens := usecases.NewTokenInvalidator(req.invalidTokenManager, req.log)

	err := invalidateTokens(ctx, user)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	return request.NewErrorResponse(models.ErrInvalidToken), nil
}
