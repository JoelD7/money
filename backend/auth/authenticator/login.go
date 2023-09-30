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
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type requestLoginHandler struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	userRepo       users.Repository
	secretsManager secrets.SecretManager
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func logInHandler(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestLoginHandler{}

	req.initLoginHandler()
	defer req.finish()

	return req.processLogin(ctx, request)
}

func (req *requestLoginHandler) initLoginHandler() {
	dynamoClient := initDynamoClient()

	req.userRepo = users.NewDynamoRepository(dynamoClient)
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLoggerWithHandler("log-in")
}

func (req *requestLoginHandler) finish() {
	defer func() {
		err := req.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLoginHandler) processLogin(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	reqBody := &Credentials{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		req.err = err
		req.log.Error("request_body_json_unmarshal_failed", err, nil)

		return apigateway.NewErrorResponse(err), nil
	}

	authenticate := usecases.NewUserAuthenticator(req.userRepo, req.log)
	generateTokens := usecases.NewTokenGenerator(req.userRepo, req.secretsManager, req.log)

	user, err := authenticate(ctx, reqBody.Username, reqBody.Password)
	if errors.Is(err, models.ErrUserNotFound) {
		return apigateway.NewErrorResponse(errUserNotFound), nil
	}

	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	accessToken, refreshToken, err := generateTokens(ctx, user)
	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	response := &accessTokenResponse{accessToken.Value}

	data, err := json.Marshal(response)
	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	setCookieHeader := map[string]string{
		"Set-Cookie": fmt.Sprintf("%s=%s; Path=/; Expires=%s; Secure; HttpOnly", refreshTokenCookieName, refreshToken.Value,
			refreshToken.Expiration.Format(time.RFC1123)),
	}

	req.log.Info("login_succeeded", nil)

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(data),
		Headers:    setCookieHeader,
	}, nil
}
