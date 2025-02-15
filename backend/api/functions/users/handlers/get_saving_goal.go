package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/usecases"
)

var (
	getSavingGoalReq  *getSavingGoalRequest
	getSavingGoalOnce sync.Once
)

type getSavingGoalRequest struct {
	startingTime   time.Time
	err            error
	savingGoalRepo savingoal.Repository
}

func (request *getSavingGoalRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	getSavingGoalOnce.Do(func() {
		request.startingTime = time.Now()
		dynamoClient := dynamo.InitClient(ctx)
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig.SavingGoalsTable)
		if err != nil {
			return
		}
	})

	return err
}

func (request *getSavingGoalRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetSavingGoalHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if getSavingGoalReq == nil {
		getSavingGoalReq = new(getSavingGoalRequest)
	}

	err := getSavingGoalReq.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_get_saving_goal_request_failed", err)

		return nil, err
	}
	defer getSavingGoalReq.finish()

	return getSavingGoalReq.process(ctx, request)
}

func (request *getSavingGoalRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingGoalID := req.PathParameters["savingGoalID"]

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	getSavingGoal := usecases.NewSavingGoalGetter(request.savingGoalRepo)

	savingGoal, err := getSavingGoal(ctx, username, savingGoalID)
	if err != nil {
		logger.Error("get_saving_goal_failed", err, req, models.Any("username", username))

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, savingGoal), nil
}
