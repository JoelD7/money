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
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     users.Repository
	incomeRepo   income.Repository
	expensesRepo expenses.Repository
}

func (request *getUserRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	guOnce.Do(func() {
		request.log = log
		request.log.SetHandler("get-user")

		dynamoClient := dynamo.InitClient(ctx)

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}

		request.incomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig.IncomeTable, envConfig.PeriodUserIncomeIndex)
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
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetUserHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if guRequest == nil {
		guRequest = new(getUserRequest)
	}

	err := guRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("get_user_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}
	defer guRequest.finish()

	return guRequest.process(ctx, req)
}

func (request *getUserRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	getUser := usecases.NewUserGetter(request.userRepo, request.incomeRepo, request.expensesRepo)

	user, err := getUser(ctx, username)
	if user != nil && user.CurrentPeriod == "" {
		request.log.Warning("user_has_no_period_set", nil, []models.LoggerObject{req})
	}

	if errors.Is(err, models.ErrIncomeNotFound) || errors.Is(err, models.ErrExpensesNotFound) {
		request.err = err
		request.log.Warning("user_remainder_could_not_be_calculated", err, []models.LoggerObject{req})

		return req.NewJSONResponse(http.StatusOK, user), nil
	}

	if errors.Is(err, models.ErrUserNotFound) {
		request.err = err
		request.log.Error("user_not_found", err, []models.LoggerObject{req})

		return req.NewErrorResponse(errors.New("user not found")), nil
	}

	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, nil)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, user), nil
}
