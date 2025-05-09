package handlers

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	guRequest *getUserRequest
	guOnce    sync.Once
)

type getUserRequest struct {
	startingTime time.Time
	err          error
	userRepo     users.Repository
	incomeRepo   income.Repository
	expensesRepo expenses.Repository
}

func (request *getUserRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	guOnce.Do(func() {
		logger.SetHandler("get-user")

		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *getUserRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetUserHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if guRequest == nil {
		guRequest = new(getUserRequest)
	}

	err := guRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("get_user_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer guRequest.finish()

	return guRequest.process(ctx, req)
}

func (request *getUserRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_user_email_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	getUser := usecases.NewUserGetter(request.userRepo, request.incomeRepo, request.expensesRepo)

	user, err := getUser(ctx, username)
	if user != nil && user.CurrentPeriod == "" {
		logger.Warning("user_has_no_period_set", nil, req)
	}

	if errors.Is(err, models.ErrIncomeNotFound) || errors.Is(err, models.ErrExpensesNotFound) {
		request.err = err
		logger.Warning("user_remainder_could_not_be_calculated", err, req)

		return req.NewJSONResponse(http.StatusOK, user), nil
	}

	if errors.Is(err, models.ErrUserNotFound) {
		request.err = err
		logger.Error("user_not_found", err, req)

		// Use a custom error instead of the error model to return 500.
		// The "username" that’s utilized to get the user from the DB is the one on the tokens that are
		// emitted by the server. Before the server creates the tokens, it supposed to have persisted the user on the DB,
		// so if the user can’t be found, then the server did something wrong.
		return req.NewErrorResponse(errors.New("user not found")), nil
	}

	if err != nil {
		request.err = err
		logger.Error("user_fetching_failed", err, nil)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, user), nil
}
