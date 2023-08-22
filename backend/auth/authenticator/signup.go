package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type requestSignUpHandler struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     *users.Repository
}

func signUpHandler(request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestSignUpHandler{
		log: logger.NewLoggerWithHandler("sign-up"),
	}

	req.initSignUpHandler()
	defer req.finish()

	return req.processSignUp(request)
}

func (req *requestSignUpHandler) initSignUpHandler() {
	dynamoClient := initDynamoClient()

	dynamoUserRepository := users.NewDynamoRepository(dynamoClient)

	req.userRepo = users.NewRepository(dynamoUserRepository)
	req.startingTime = time.Now()
	req.log = logger.NewLogger()
}

func (req *requestSignUpHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestSignUpHandler) processSignUp(request *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	reqBody := new(signUpBody)

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		req.err = err
		req.log.Error("request_body_json_unmarshal_failed", err, nil)

		return getErrorResponse(err)
	}

	saveNewUser := usecases.NewUserCreator(req.userRepo, req.log)

	err = saveNewUser(ctx, reqBody.FullName, reqBody.Email, reqBody.Password)
	if err != nil {
		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusCreated,
	}, nil
}
