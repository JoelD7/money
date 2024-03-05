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

func logInHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestLoginHandler{}

	req.initLoginHandler(log)
	defer req.finish()

	return req.processLogin(ctx, request)
}

func (req *requestLoginHandler) initLoginHandler(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	req.userRepo = users.NewDynamoRepository(dynamoClient)
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = log
	req.log.SetHandler("login")
}

func (req *requestLoginHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLoginHandler) processLogin(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	reqBody, err := validateLoginInput(request)
	if err != nil {
		req.err = err
		req.log.Error("validate_input_failed", err, []models.LoggerObject{request})

		return request.NewErrorResponse(err), nil
	}

	authenticate := usecases.NewUserAuthenticator(req.userRepo, req.log)
	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager, req.log)

	user, err := authenticate(ctx, reqBody.Username, reqBody.Password)
	if errors.Is(err, models.ErrUserNotFound) {
		return request.NewErrorResponse(errUserNotFound), nil
	}

	if err != nil {
		return request.NewErrorResponse(err), nil
	}

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

	req.log.Info("login_succeeded", nil)

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(data),
		Headers:    setCookieHeader,
	}, nil
}

func validateLoginInput(request *apigateway.Request) (*Credentials, error) {
	reqBody := new(Credentials)

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return nil, err
	}

	err = validateCredentials(reqBody.Username, reqBody.Password)
	if err != nil {
		return nil, err
	}

	return reqBody, nil
}
