package main

import (
	"context"
	"encoding/json"
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

var signUpRequest *requestSignUpHandler
var signUpOnce sync.Once

type requestSignUpHandler struct {
	startingTime   time.Time
	err            error
	userRepo       users.Repository
	secretsManager secrets.SecretManager
}

func signUpHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if signUpRequest == nil {
		signUpRequest = new(requestSignUpHandler)
	}

	err := signUpRequest.initSignUpHandler(ctx, envConfig)
	if err != nil {
		signUpRequest.err = err

		logger.Error("sign_up_init_failed", err, request)

		return request.NewErrorResponse(err), nil
	}
	defer signUpRequest.finish()

	return signUpRequest.processSignUp(ctx, request)
}

func (req *requestSignUpHandler) initSignUpHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	signUpOnce.Do(func() {
		logger.SetHandler("sign-up")
		dynamoClient := dynamo.InitClient(ctx)

		req.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
		req.secretsManager = secrets.NewAWSSecretManager()
	})
	req.startingTime = time.Now()

	return err
}

func (req *requestSignUpHandler) finish() {
	logger.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestSignUpHandler) processSignUp(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	reqBody, err := validateSingUpInput(request)
	if err != nil {
		req.err = err
		logger.Error("validate_input_failed", err, request)

		return request.NewErrorResponse(err), nil
	}

	saveNewUser := usecases.NewUserCreator(req.userRepo)

	newUser, err := saveNewUser(ctx, reqBody.FullName, reqBody.Username, reqBody.Password)
	if err != nil {
		req.err = err
		logger.Error("save_new_user_failed", err, request)

		return request.NewErrorResponse(err), nil
	}

	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager)

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

	logger.Info("signup_succeeded", request)

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
