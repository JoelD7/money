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
	"github.com/JoelD7/money/backend/storage/dynamo"
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

	startingTime        time.Time
	err                 error
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
	secretsManager      secrets.SecretManager
}

func tokenHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if req == nil {
		req = new(requestTokenHandler)
	}

	err := req.initTokenHandler(ctx, envConfig)
	if err != nil {
		req.err = err
		logger.Error("token_init_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}
	defer req.finish()

	return req.processToken(ctx, request)
}

func (req *requestTokenHandler) initTokenHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	once.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		req.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
		req.invalidTokenManager = cache.NewRedisCache()
		req.secretsManager = secrets.NewAWSSecretManager()
		logger.SetHandler("token")
	})

	req.startingTime = time.Now()
	req.err = nil

	return err
}

func (req *requestTokenHandler) finish() {
	logger.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestTokenHandler) processToken(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	var err error
	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		logger.Error("getting_refresh_token_cookie_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	validateRefreshToken := usecases.NewRefreshTokenValidator(req.userRepo)

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

	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager)

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

	logger.Info("new_tokens_issued_successfully", user)

	return request.NewJSONResponse(http.StatusOK, string(data), apigateway.Header{Key: "Set-Cookie", Value: cookieStr}), nil
}

func (req *requestTokenHandler) handleValidationError(ctx context.Context, user *models.User, request *apigateway.Request) (*apigateway.Response, error) {
	invalidateTokens := usecases.NewTokenInvalidator(req.invalidTokenManager)

	err := invalidateTokens(ctx, user)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	return request.NewErrorResponse(models.ErrInvalidToken), nil
}
