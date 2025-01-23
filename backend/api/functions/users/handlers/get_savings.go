package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gssRequest *getSavingsRequest
var gssOnce sync.Once

type getSavingsRequest struct {
	username       string
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
	*models.QueryParameters
}

type savingsResponse struct {
	Savings []*models.Saving `json:"savings"`
	NextKey string           `json:"next_key"`
}

func (request *getSavingsRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gssOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig.SavingsTable, envConfig.PeriodSavingIndexName, envConfig.SavingGoalSavingIndexName)
		if err != nil {
			return
		}
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig.SavingGoalsTable)
		if err != nil {
			return
		}
		logger.SetHandler("get-savings")
	})
	request.startingTime = time.Now()

	return err
}

func (request *getSavingsRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetSavingsHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gssRequest == nil {
		gssRequest = new(getSavingsRequest)
	}

	err := gssRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("get_savings_request_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer gssRequest.finish()

	err = gssRequest.prepareRequest(req)
	if err != nil {
		return req.NewErrorResponse(err), nil
	}

	return gssRequest.routeToHandlers(ctx, req)
}

func (request *getSavingsRequest) prepareRequest(req *apigateway.Request) error {
	var err error

	request.username, err = apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_user_email_from_context_failed", err, req)

		return err
	}

	err = validate.Email(request.username)
	if err != nil {
		logger.Error("invalid_username", err, logger.MapToLoggerObject("user_data", map[string]interface{}{
			"s_username": request.username,
		}))

		return err
	}

	request.QueryParameters, err = req.GetQueryParameters()
	if err != nil {
		logger.Error("get_request_params_failed", err, req)

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
	getSavings := usecases.NewSavingsGetter(request.savingsRepo, request.savingGoalRepo)

	userSavings, nextKey, err := getSavings(ctx, request.username, request.StartKey, request.PageSize)
	if err != nil {
		logger.Error("savings_fetch_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		logger.Error("savings_marshal_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period := req.QueryStringParameters["period"]

	getSavingsByPeriod := usecases.NewSavingByPeriodGetter(request.savingsRepo, request.savingGoalRepo)

	userSavings, nextKey, err := getSavingsByPeriod(ctx, request.username, request.StartKey, period, request.PageSize)
	if err != nil {
		logger.Error("savings_fetch_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		logger.Error("savings_marshal_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsBySavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingGoal := req.QueryStringParameters["saving_goal_id"]

	getSavingsBySavingGoal := usecases.NewSavingBySavingGoalGetter(request.savingsRepo, request.savingGoalRepo)

	userSavings, nextKey, err := getSavingsBySavingGoal(ctx, request.StartKey, savingGoal, request.PageSize)
	if err != nil {
		logger.Error("savings_fetch_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		logger.Error("savings_marshal_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       responseJSON,
	}, nil
}

func (request *getSavingsRequest) getUserSavingsByPeriodAndSavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	period := req.QueryStringParameters["period"]
	savingGoal := req.QueryStringParameters["saving_goal_id"]

	getSavingsBySavingGoalAndPeriod := usecases.NewSavingBySavingGoalAndPeriodGetter(request.savingsRepo, request.savingGoalRepo)

	userSavings, nextKey, err := getSavingsBySavingGoalAndPeriod(ctx, request.StartKey, savingGoal, period, request.PageSize)
	if err != nil {
		logger.Error("savings_fetch_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	responseJSON, err := request.getSavingsResponse(userSavings, nextKey)
	if err != nil {
		logger.Error("savings_marshal_failed", err, req)

		return req.NewErrorResponse(err), nil
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
