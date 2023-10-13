package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"strconv"
	"time"
)

type expensesResponse struct {
	Expenses []*models.Expense `json:"expenses"`
	NextKey  string            `json:"next_key"`
}

type getExpensesRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	expensesRepo expenses.Repository
	userRepo     users.Repository
	username     string
	startKey     string
	pageSize     int
}

func (request *getExpensesRequest) init() {
	dynamoClient := initDynamoClient()

	request.expensesRepo = expenses.NewDynamoRepository(dynamoClient)
	request.userRepo = users.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getExpensesRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getExpensesHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getExpensesRequest)

	request.init()
	defer request.finish()

	return request.routeToHandlers(ctx, req)
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

	request.startKey, request.pageSize, err = getRequestParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return err
	}

	return nil
}

func getRequestParams(req *apigateway.Request) (string, int, error) {
	pageSizeParam := 0
	var err error

	if req.QueryStringParameters["page_size"] != "" {
		pageSizeParam, err = strconv.Atoi(req.QueryStringParameters["page_size"])
		if err != nil {
			return "", 0, err
		}
	}

	return req.QueryStringParameters["start_key"], pageSizeParam, nil
}

func (request *getExpensesRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	_, categoryOk := req.MultiValueQueryStringParameters["category"]
	_, periodOk := req.QueryStringParameters["period"]

	if categoryOk && !periodOk {
		return request.getByCategories(ctx, req)
	}

	if !categoryOk && periodOk {
		return request.getByPeriod(ctx, req)
	}

	if categoryOk && periodOk {
		return request.getByCategoriesAndPeriod(ctx, req)
	}

	return request.getAll(ctx, req)
}

func (request *getExpensesRequest) getByCategories(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	categories, _ := req.MultiValueQueryStringParameters["category"]

	getExpensesByCategory := usecases.NewExpensesByCategoriesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByCategory(ctx, request.username, request.startKey, categories, request.pageSize)
	if errors.Is(err, models.ErrCategoryNameSettingFailed) {
		request.log.Error("set_expense_category_name_failed", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
			Expenses: userExpenses,
			NextKey:  nextKey,
		}), nil
	}

	if err != nil {
		request.log.Error("get_expenses_by_category_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *getExpensesRequest) getByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period, _ := req.QueryStringParameters["period"]

	getExpensesByPeriod := usecases.NewExpensesByPeriodGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByPeriod(ctx, request.username, period, request.startKey, request.pageSize)
	if errors.Is(err, models.ErrCategoryNameSettingFailed) {
		request.log.Error("set_expense_category_name_failed", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
			Expenses: userExpenses,
			NextKey:  nextKey,
		}), nil
	}

	if err != nil {
		request.log.Error("get_expenses_by_period_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return nil, nil
}

func (request *getExpensesRequest) getByCategoriesAndPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	categories, _ := req.MultiValueQueryStringParameters["category"]
	period, _ := req.QueryStringParameters["period"]

	getExpensesByPeriodAndCategories := usecases.NewExpensesByPeriodAndCategoriesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpensesByPeriodAndCategories(ctx, request.username, period, request.startKey, categories, request.pageSize)
	if errors.Is(err, models.ErrCategoryNameSettingFailed) {
		request.log.Error("set_expense_category_name_failed", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
			Expenses: userExpenses,
			NextKey:  nextKey,
		}), nil
	}

	if err != nil {
		request.log.Error("get_expenses_by_period_and_categories_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}

func (request *getExpensesRequest) getAll(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getExpenses := usecases.NewExpensesGetter(request.expensesRepo, request.userRepo)

	userExpenses, nextKey, err := getExpenses(ctx, request.username, request.startKey, request.pageSize)
	if errors.Is(err, models.ErrCategoryNameSettingFailed) {
		request.log.Error("set_expense_category_name_failed", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
			Expenses: userExpenses,
			NextKey:  nextKey,
		}), nil
	}

	if err != nil {
		request.log.Error("get_expenses_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, &expensesResponse{
		Expenses: userExpenses,
		NextKey:  nextKey,
	}), nil
}
