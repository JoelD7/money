package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	userRepo       *users.Repository
	secretsManager secrets.SecretManager
}

func logInHandler(request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestLoginHandler{
		log: logger.NewLoggerWithHandler("log-in"),
	}

	req.initLoginHandler()
	defer req.finish()

	return req.processLogin(request)
}

func (req *requestLoginHandler) initLoginHandler() {
	dynamoClient := initDynamoClient()

	dynamoUserRepository := users.NewDynamoRepository(dynamoClient)

	req.userRepo = users.NewRepository(dynamoUserRepository)
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLogger()
}

func (req *requestLoginHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestLoginHandler) processLogin(request *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	reqBody := &Credentials{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		req.err = err
		req.log.Error("request_body_json_unmarshal_failed", err, nil)

		return getErrorResponse(err)
	}

	authenticate := usecases.NewUserAuthenticator(req.userRepo, req.log)
	generateTokens := usecases.NewUserTokenGenerator(req.userRepo, req.secretsManager, req.log)

	user, err := authenticate(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		return getErrorResponse(err)
	}

	accessToken, refreshToken, err := generateTokens(ctx, user)
	if err != nil {
		return getErrorResponse(err)
	}

	setCookieHeader := map[string][]string{
		"Set-Cookie": {
			fmt.Sprintf("%s=%s; Path=/; Expires=%s; Secure; HttpOnly", accessTokenCookieName, accessToken.Value,
				accessToken.Expiration.Format(time.RFC1123)),
			fmt.Sprintf("%s=%s; Path=/; Expires=%s; Secure; HttpOnly", refreshTokenCookieName, refreshToken.Value,
				refreshToken.Expiration.Format(time.RFC1123)),
		},
	}

	req.log.Info("login_succeeded", nil)

	return &apigateway.Response{
		StatusCode:        http.StatusOK,
		MultiValueHeaders: setCookieHeader,
	}, nil
}
