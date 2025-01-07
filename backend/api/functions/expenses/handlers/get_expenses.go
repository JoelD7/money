package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gesExpensesRequest *GetExpensesRequest
var gesOnce sync.Once

type ExpensesResponse struct {
	Expenses []*models.Expense `json:"expenses"`
	NextKey  string            `json:"next_key"`
}

type GetExpensesRequest struct {
	Username string
	*models.QueryParameters

	Log          logger.LogAPI
	ExpensesRepo expenses.Repository
	UserRepo     users.Repository

	startingTime time.Time
	err          error
}

func (request *GetExpensesRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gesOnce.Do(func() {
		request.Log = log
		dynamoClient := dynamo.InitClient(ctx)

		request.ExpensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.UserRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *GetExpensesRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetExpenses(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gesExpensesRequest == nil {
		gesExpensesRequest = new(GetExpensesRequest)
	}

	err := gesExpensesRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("get_expenses_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer gesExpensesRequest.finish()

	err = gesExpensesRequest.prepareRequest(req)
	if err != nil {
		return req.NewErrorResponse(err), nil
	}

	return gesExpensesRequest.routeToHandlers(ctx, req)
}

func (request *GetExpensesRequest) prepareRequest(req *apigateway.Request) error {
	var err error

	request.Username, err = apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.Log.Error("get_user_email_from_context_failed", err, req)

		return err
	}

	err = validate.Email(request.Username)
	if err != nil {
		request.Log.Error("invalid_username", err,
			models.Any("user_data", map[string]interface{}{
				"s_username": request.Username,
			}),
		)

		return err
	}

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		request.Log.Error("get_request_params_failed", err, req)

		return err
	}

	err = validate.SortBy(request.SortBy)
	if err != nil {
		return err
	}

	err = validate.SortType(request.SortType)
	if err != nil {
		return err
	}

	return nil
}

func (request *GetExpensesRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	if len(request.Categories) > 0 && request.Period == "" {
		return request.getByCategories(ctx, req)
	}

	if len(request.Categories) == 0 && request.Period != "" {
		return request.GetByPeriod(ctx, req)
	}

	if len(request.Categories) > 0 && request.Period != "" {
		return request.getByCategoriesAndPeriod(ctx, req)
	}

	return request.getAll(ctx, req)
}

func (request *GetExpensesRequest) getByCategories(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByCategory := usecases.NewExpensesByCategoriesGetter(request.ExpensesRepo, request.UserRepo)

	userExpenses, nextKey, err := getExpensesByCategory(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.Log.Error("get_expenses_by_category_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &ExpensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *GetExpensesRequest) GetByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByPeriod := usecases.NewExpensesByPeriodGetter(request.ExpensesRepo, request.UserRepo)

	userExpenses, nextKey, err := getExpensesByPeriod(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.Log.Error("get_expenses_by_period_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &ExpensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *GetExpensesRequest) getByCategoriesAndPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByPeriodAndCategories := usecases.NewExpensesByPeriodAndCategoriesGetter(request.ExpensesRepo, request.UserRepo)

	userExpenses, nextKey, err := getExpensesByPeriodAndCategories(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.Log.Error("get_expenses_by_period_and_categories_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &ExpensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *GetExpensesRequest) getAll(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpenses := usecases.NewExpensesGetter(request.ExpensesRepo, request.UserRepo)

	userExpenses, nextKey, err := getExpenses(ctx, request.Username, request.QueryParameters)
	if err != nil {
		request.Log.Error("get_expenses_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &ExpensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}
