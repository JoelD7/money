package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type getUserRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
	incomeRepo   income.Repository
	expensesRepo expenses.Repository
}

func (request *getUserRequest) init() {
	dynamoClient := initDynamoClient()

	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.incomeRepo = income.NewDynamoRepository(dynamoClient)
	request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLoggerWithHandler("get-user")
}

func (request *getUserRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getUserHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getUserRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getUserRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := getUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	getUser := usecases.NewUserGetter(request.userRepo, request.incomeRepo, request.expensesRepo)

	user, err := getUser(ctx, username)
	if user != nil && user.CurrentPeriod == "" {
		request.log.Warning("user_has_no_period_set", nil, []models.LoggerObject{req})
	}

	if errors.Is(err, models.ErrIncomeNotFound) || errors.Is(err, models.ErrExpensesNotFound) {
		request.err = err
		request.log.Warning("user_remainder_could_not_be_calculated", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, user), nil
	}

	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, nil)

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, user), nil
}
