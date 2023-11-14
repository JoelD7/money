package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type getSavingsRequest struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
	username       string
	startKey       string
	pageSize       int
}

type savingsResponse struct {
	Savings []*models.Saving `json:"savings"`
	NextKey string           `json:"next_key"`
}

func (request *getSavingsRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.savingGoalRepo = savingoal.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getSavingsRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getSavingsHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getSavingsRequest)

	request.init()
	defer request.finish()

	err := request.prepareRequest(req)
	if err != nil {
		return apigateway.NewErrorResponse(err), nil
	}

	return request.routeToHandlers(ctx, req)
}

func (request *getSavingsRequest) prepareRequest(req *apigateway.Request) error {
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

	request.startKey, request.pageSize, err = getRequestQueryParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return err
	}

	return nil
}

func (request *getSavingsRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	_, periodOk := req.QueryStringParameters["period"]
	_, savingGoalOk := req.QueryStringParameters["saving_goal_id"]

	if periodOk && !savingGoalOk {
		return request.getUserSavingsByPeriod(ctx, req)
	}

	if !periodOk && savingGoalOk {
		return request.getUserSavingsBySavingGoal(ctx, req)
	}

	if periodOk && savingGoalOk {
		return request.getUserSavingsByPeriodAndSavingGoal(ctx, req)
	}

	return request.getUserSavings(ctx, req)
}

func (request *getSavingsRequest) getUserSavings(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getSavings := usecases.NewSavingsGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	userSavings, nextKey, err := getSavings(ctx, request.username, request.startKey, request.pageSize)
	if err != nil {
		request.log.Error("savings_fetch_failed", err, []models.LoggerObject{
			req,
		})

		return apigateway.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		request.log.Error("savings_marshal_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period := req.QueryStringParameters["period"]

	getSavingsByPeriod := usecases.NewSavingByPeriodGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	userSavings, nextKey, err := getSavingsByPeriod(ctx, request.username, request.startKey, period, request.pageSize)
	if err != nil {
		request.log.Error("savings_fetch_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		request.log.Error("savings_marshal_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsBySavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingGoal := req.QueryStringParameters["saving_goal_id"]

	getSavingsBySavingGoal := usecases.NewSavingBySavingGoalGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	userSavings, nextKey, err := getSavingsBySavingGoal(ctx, request.startKey, savingGoal, request.pageSize)
	if err != nil {
		request.log.Error("savings_fetch_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		request.log.Error("savings_marshal_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsByPeriodAndSavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period := req.QueryStringParameters["period"]
	savingGoal := req.QueryStringParameters["saving_goal_id"]

	getSavingsBySavingGoalAndPeriod := usecases.NewSavingBySavingGoalAndPeriodGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	userSavings, nextKey, err := getSavingsBySavingGoalAndPeriod(ctx, request.startKey, savingGoal, period, request.pageSize)
	if err != nil {
		request.log.Error("savings_fetch_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		request.log.Error("savings_marshal_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getSavingsResponse(savings []*models.Saving, nextKey string) (string, error) {
	response := &savingsResponse{savings, nextKey}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(responseJSON), nil
}
