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
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var loginRequest *requestLoginHandler
var loginOnce sync.Once

type requestLoginHandler struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	userRepo       users.Repository
	secretsManager secrets.SecretManager
}

type accessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func logInHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	if loginRequest == nil {
		loginRequest = new(requestLoginHandler)
	}

	loginRequest.initLoginHandler(ctx, log)
	defer loginRequest.finish()

	return loginRequest.processLogin(ctx, request)
}

func (req *requestLoginHandler) initLoginHandler(ctx context.Context, log logger.LogAPI) {
	loginOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		req.userRepo = users.NewDynamoRepository(dynamoClient)
		req.secretsManager = secrets.NewAWSSecretManager()
		req.log = log
		req.log.SetHandler("login")
	})
	req.startingTime = time.Now()
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

	cookieStr := getRefreshTokenCookieStr(refreshToken.Value, refreshToken.Expiration)

	req.log.Info("login_succeeded", nil)

	return request.NewJSONResponse(http.StatusOK, string(data), apigateway.Header{Key: "Set-Cookie", Value: cookieStr}), nil
}

func getRefreshTokenCookieStr(value string, expiration time.Time) string {
	return fmt.Sprintf("%s=%s; Expires=%s; Path=/; Secure; SameSite=None; HttpOnly;", refreshTokenCookieName, value,
		expiration.Format(time.RFC1123))
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
