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

var gesExpensesRequest *getExpensesRequest
var gesOnce sync.Once

type expensesResponse struct {
	Expenses []*models.Expense `json:"expenses"`
	NextKey  string            `json:"next_key"`
}

type getExpensesRequest struct {
	username string
	apigateway.QueryParameters

	log          logger.LogAPI
	expensesRepo expenses.Repository
	userRepo     users.Repository

	startingTime time.Time
	err          error
}

func (request *getExpensesRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gesOnce.Do(func() {
		request.log = log
		dynamoClient := dynamo.InitClient(ctx)

		request.expensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig.ExpensesTable, envConfig.ExpensesRecurringTable, envConfig.PeriodUserExpenseIndex)
		if err != nil {
			return
		}

		request.userRepo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *getExpensesRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetExpenses(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gesExpensesRequest == nil {
		gesExpensesRequest = new(getExpensesRequest)
	}

	err := gesExpensesRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("get_expenses_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}
	defer gesExpensesRequest.finish()

	err = gesExpensesRequest.prepareRequest(req)
	if err != nil {
		return req.NewErrorResponse(err), nil
	}

	return gesExpensesRequest.routeToHandlers(ctx, req)
}

func (request *getExpensesRequest) prepareRequest(req *apigateway.Request) error {
	var err error

	request.username, err = apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return err
	}

	err = validate.Email(request.username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{
			request.log.MapToLoggerObject("user_data", map[string]interface{}{
				"s_username": request.username,
			}),
		})

		return err
	}

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return err
	}

	return nil
}

func (request *getExpensesRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	if len(request.Categories) > 0 && request.Period == "" {
		return request.getByCategories(ctx, req)
	}

	if len(request.Categories) == 0 && request.Period != "" {
		return request.getByPeriod(ctx, req)
	}

	if len(request.Categories) > 0 && request.Period != "" {
		return request.getByCategoriesAndPeriod(ctx, req)
	}

	return request.getAll(ctx, req)
}

func (request *getExpensesRequest) getByCategories(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByCategory := usecases.NewExpensesByCategoriesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByCategory(ctx, request.username, request.StartKey, request.Categories, request.PageSize)
	if err != nil {
		request.log.Error("get_expenses_by_category_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *getExpensesRequest) getByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByPeriod := usecases.NewExpensesByPeriodGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByPeriod(ctx, request.username, request.Period, request.StartKey, request.PageSize)
	if err != nil {
		request.log.Error("get_expenses_by_period_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *getExpensesRequest) getByCategoriesAndPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpensesByPeriodAndCategories := usecases.NewExpensesByPeriodAndCategoriesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByPeriodAndCategories(ctx, request.username, request.Period, request.StartKey, request.Categories, request.PageSize)
	if err != nil {
		request.log.Error("get_expenses_by_period_and_categories_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *getExpensesRequest) getAll(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpenses := usecases.NewExpensesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpenses(ctx, request.username, request.StartKey, request.PageSize)
	if err != nil {
		request.log.Error("get_expenses_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}
