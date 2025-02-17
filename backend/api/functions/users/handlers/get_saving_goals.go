package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	getSavingGoalsReq  *getSavingGoalsRequest
	getSavingGoalsOnce sync.Once
)

type getSavingGoalsRequest struct {
	startingTime   time.Time
	err            error
	savingGoalRepo savingoal.Repository
	queryParams    *models.QueryParameters
}

type SavingGoalsResponse struct {
	SavingGoals []*models.SavingGoal `json:"saving_goals"`
	NextKey     string               `json:"next_key"`
}

func (request *getSavingGoalsRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	getSavingGoalsOnce.Do(func() {
		request.startingTime = time.Now()
		dynamoClient := dynamo.InitClient(ctx)
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig.SavingGoalsTable)
		if err != nil {
			return
		}
	})

	return err
}

func (request *getSavingGoalsRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetSavingGoalsHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if getSavingGoalsReq == nil {
		getSavingGoalsReq = new(getSavingGoalsRequest)
	}

	err := getSavingGoalsReq.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_get_saving_goals_request_failed", err)

		return nil, err
	}
	defer getSavingGoalsReq.finish()

	getSavingGoalsReq.queryParams, err = request.GetQueryParameters()
	if err != nil {
		logger.Error("get_query_parameters_failed", err, request)

		return request.NewErrorResponse(err), nil
	}

	err = getSavingGoalsReq.validateQueryParams()
	if err != nil {
		logger.Error("validate_query_params_failed", err, request)

		return request.NewErrorResponse(err), nil
	}

	return getSavingGoalsReq.process(ctx, request)
}

func (request *getSavingGoalsRequest) validateQueryParams() error {
	err := validate.SortBy(request.queryParams.SortBy, validate.SortByModelSavingGoals)
	if err != nil {
		return err
	}

	err = validate.SortType(request.queryParams.SortType)
	if err != nil {
		return err
	}

	return nil
}

func (request *getSavingGoalsRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	getSavingGoals := usecases.NewSavingGoalsGetter(request.savingGoalRepo)

	savingGoals, nextKey, err := getSavingGoals(ctx, username, request.queryParams)
	if err != nil {
		logger.Error("get_saving_goals_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, &SavingGoalsResponse{
		SavingGoals: savingGoals,
		NextKey:     nextKey,
	}), nil
}
