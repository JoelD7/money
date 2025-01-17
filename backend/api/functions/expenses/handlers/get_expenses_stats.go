package handlers

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	gesExpenseRequest *GetExpensesStatsRequest
	gestExpenseOnce   sync.Once
)

type GetExpensesStatsRequest struct {
	startingTime time.Time
	err          error
	ExpensesRepo expenses.Repository
	UserRepo     users.Repository
}

func (request *GetExpensesStatsRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gestExpenseOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.ExpensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *GetExpensesStatsRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetExpensesStats(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gesExpenseRequest == nil {
		gesExpenseRequest = new(GetExpensesStatsRequest)
	}

	err := gesExpenseRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("get_expenses_stats_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	defer gesExpenseRequest.finish()

	return gesExpenseRequest.Process(ctx, req)
}

func (request *GetExpensesStatsRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok {
		logger.Error("missing_period_id", fmt.Errorf("period ID not in path parameters"), req)

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	getCategoryExpensesSummary := usecases.NewCategoryExpenseSummaryGetter(request.ExpensesRepo)
	categoryExpenseSummary, err := getCategoryExpensesSummary(ctx, username, periodID)
	if err != nil {
		logger.Error("get_expenses_stats_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, categoryExpenseSummary), nil
}
