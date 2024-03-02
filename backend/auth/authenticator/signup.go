package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
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
	userRepo     users.Repository
}

func signUpHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestSignUpHandler{}

	req.initSignUpHandler(log)
	defer req.finish()

	return req.processSignUp(ctx, request)
}

func (req *requestSignUpHandler) initSignUpHandler(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	req.userRepo = users.NewDynamoRepository(dynamoClient)
	req.startingTime = time.Now()
	req.log = log
	req.log.SetHandler("sign-up")
}

func (req *requestSignUpHandler) finish() {
	defer func() {
		err := req.log.Close()
		if err != nil {
			panic(err)
		}
	}()

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

	return request.NewJSONResponse(http.StatusCreated, nil), nil
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
