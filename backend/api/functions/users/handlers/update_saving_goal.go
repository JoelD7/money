package handlers

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	updateSavingGoalsReq  *updateSavingGoalsRequest
	updateSavingGoalsOnce sync.Once
)

type updateSavingGoalsRequest struct {
	startingTime   time.Time
	err            error
	savingGoalRepo savingoal.Repository
}

func (request *updateSavingGoalsRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	updateSavingGoalsOnce.Do(func() {
		request.startingTime = time.Now()
		dynamoClient := dynamo.InitClient(ctx)
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})

	return err
}

func (request *updateSavingGoalsRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func UpdateSavingGoalsHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if updateSavingGoalsReq == nil {
		updateSavingGoalsReq = new(updateSavingGoalsRequest)
	}

	err := updateSavingGoalsReq.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_update_saving_goals_request_failed", err)

		return request.NewErrorResponse(err), nil
	}
	defer updateSavingGoalsReq.finish()

	return updateSavingGoalsReq.process(ctx, request)
}

func (request *updateSavingGoalsRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	savingGoalID, ok := req.PathParameters["savingGoalID"]
	if !ok {
		err := fmt.Errorf("missing saving goal ID")
		logger.Error("missing_saving_goal_id", err, req)

		return req.NewErrorResponse(err), nil
	}

	savingGoal, err := validateSavingGoalBody(req)
	if err != nil {
		logger.Error("validate_saving_goal_body_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	updateSavingGoal := usecases.NewSavingGoalUpdator(request.savingGoalRepo)

	updatedGoal, err := updateSavingGoal(ctx, username, savingGoalID, savingGoal)
	if err != nil {
		logger.Error("update_saving_goal_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, updatedGoal), nil
}
