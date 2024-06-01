package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var signUpRequest *requestSignUpHandler
var signUpOnce sync.Once

type requestSignUpHandler struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	userRepo       users.Repository
	secretsManager secrets.SecretManager
}

func signUpHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	if signUpRequest == nil {
		signUpRequest = new(requestSignUpHandler)
	}

	signUpRequest.initSignUpHandler(log)
	defer signUpRequest.finish()

	return signUpRequest.processSignUp(ctx, request)
}

func (req *requestSignUpHandler) initSignUpHandler(log logger.LogAPI) {
	signUpOnce.Do(func() {
		dynamoClient := initDynamoClient()

		req.userRepo = users.NewDynamoRepository(dynamoClient)
		req.log = log
		req.log.SetHandler("sign-up")
		req.secretsManager = secrets.NewAWSSecretManager()
	})
	req.startingTime = time.Now()
}

func (req *requestSignUpHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestSignUpHandler) processSignUp(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	reqBody, err := validateSingUpInput(request)
	if err != nil {
		req.err = err
		req.log.Error("validate_input_failed", err, []models.LoggerObject{request})

		return request.NewErrorResponse(err), nil
	}

	saveNewUser := usecases.NewUserCreator(req.userRepo, req.log)

	err = saveNewUser(ctx, reqBody.FullName, reqBody.Username, reqBody.Password)
	if err != nil {
		req.err = err
		req.log.Error("save_new_user_failed", err, []models.LoggerObject{request})

		return request.NewErrorResponse(err), nil
	}

	newUser := &models.User{
		Username: reqBody.Username,
		FullName: reqBody.FullName,
	}

	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager, req.log)

	accessToken, refreshToken, err := generateTokens(ctx, newUser)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	response := &accessTokenResponse{accessToken.Value}

	data, err := json.Marshal(response)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	cookieStr := getRefreshTokenCookieStr(refreshToken.Value, refreshToken.Expiration)

	return request.NewJSONResponse(http.StatusCreated, string(data), apigateway.Header{Key: "Set-Cookie", Value: cookieStr}), nil
}

func validateSingUpInput(request *apigateway.Request) (*signUpBody, error) {
	reqBody := new(signUpBody)

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
