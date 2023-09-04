package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type requestTokenHandler struct {
	RefreshToken string `json:"refresh_token,omitempty"`

	log                 logger.LogAPI
	startingTime        time.Time
	err                 error
	userRepo            *users.Repository
	invalidTokenManager cache.InvalidTokenManager
	secretsManager      secrets.SecretManager
}

func tokenHandler(request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestTokenHandler{}

	req.initTokenHandler()
	defer req.finish()

	return req.processToken(request)
}

func (req *requestTokenHandler) initTokenHandler() {
	dynamoClient := initDynamoClient()

	dynamoUserRepository := users.NewDynamoRepository(dynamoClient)
	req.userRepo = users.NewRepository(dynamoUserRepository)

	req.invalidTokenManager = cache.NewRedisCache()
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLoggerWithHandler("token")
}

func (req *requestTokenHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestTokenHandler) processToken(request *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	var err error
	req.RefreshToken, err = getRefreshTokenCookie(request)
	if err != nil {
		req.err = err
		req.log.Error("getting_refresh_token_cookie_failed", err, nil)

		return getErrorResponse(err)
	}

	validateRefreshToken := usecases.NewRefreshTokenValidator(req.userRepo, req.log)

	user, err := validateRefreshToken(ctx, req.RefreshToken)
	if err != nil && errors.Is(err, models.ErrInvalidToken) {
		req.err = err

		return req.handleValidationError(ctx, user)
	}

	if err != nil {
		req.err = err

		return getErrorResponse(err)
	}

	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager, req.log)

	accessToken, refreshToken, err := generateTokens(ctx, user)
	if err != nil {
		return getErrorResponse(err)
	}

	response := &accessTokenResponse{accessToken.Value}

	data, err := json.Marshal(response)
	if err != nil {
		return getErrorResponse(err)
	}

	setCookieHeader := map[string]string{
		"Set-Cookie": fmt.Sprintf("%s=%s; Path=/; Expires=%s; Secure; HttpOnly", refreshTokenCookieName, refreshToken.Value,
			refreshToken.Expiration.Format(time.RFC1123)),
	}

	req.log.Info("new_tokens_issued_successfully", []models.LoggerObject{user})

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(data),
		Headers:    setCookieHeader,
	}, nil
}

func (req *requestTokenHandler) handleValidationError(ctx context.Context, user *models.User) (*apigateway.Response, error) {
	invalidateTokens := usecases.NewTokenInvalidator(req.invalidTokenManager, req.log)

	err := invalidateTokens(ctx, user)
	if err != nil {
		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       models.ErrInvalidToken.Error(),
	}, nil
}
