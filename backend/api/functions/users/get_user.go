package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
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
	request.log = logger.NewLogger()
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
	userID := req.PathParameters["user-id"]

	getUser := usecases.NewUserGetter(request.userRepo, request.incomeRepo, request.expensesRepo)

	user, err := getUser(ctx, userID)
	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, nil)

		return apigateway.NewErrorResponse(errUserFetchingFailed), nil
	}

	if user == nil {
		request.err = err
		request.log.Error("user_not_found", err, nil)

		return apigateway.NewErrorResponse(ErrNotFound), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, user), nil
}
