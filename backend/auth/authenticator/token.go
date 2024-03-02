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
	userRepo            users.Repository
	invalidTokenManager cache.InvalidTokenManager
	secretsManager      secrets.SecretManager
}

func tokenHandler(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestTokenHandler{}

	req.initTokenHandler()
	defer req.finish()

	return req.processToken(ctx, request)
}

func (req *requestTokenHandler) initTokenHandler() {
	dynamoClient := initDynamoClient()

	req.userRepo = users.NewDynamoRepository(dynamoClient)

	req.invalidTokenManager = cache.NewRedisCache()
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLoggerWithHandler("token")
}

func (req *requestTokenHandler) finish() {
	defer func() {
		err := req.log.Close()
		if err != nil {
			panic(err)
		}
	}()

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

	setCookieHeader := map[string]string{
		"Set-Cookie": fmt.Sprintf("%s=%s; Path=/; Expires=%s; Secure; HttpOnly", refreshTokenCookieName, refreshToken.Value,
			refreshToken.Expiration.Format(time.RFC1123)),
	}

	fmt.Println("hi")
	req.log.Info("new_tokens_issued_successfully", []models.LoggerObject{user})

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(data),
		Headers:    setCookieHeader,
	}, nil
}

func (req *requestTokenHandler) handleValidationError(ctx context.Context, user *models.User, request *apigateway.Request) (*apigateway.Response, error) {
	invalidateTokens := usecases.NewTokenInvalidator(req.invalidTokenManager, req.log)

	err := invalidateTokens(ctx, user)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       models.ErrInvalidToken.Error(),
	}, nil
}
