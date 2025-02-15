package handlers

import (
	"context"
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
	dsgRequest *deleteSavingGoalRequest
	dsg        sync.Once
)

type deleteSavingGoalRequest struct {
	startingTime   time.Time
	err            error
	savingGoalRepo savingoal.Repository
}

func (request *deleteSavingGoalRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	dsg.Do(func() {
		request.startingTime = time.Now()
		dynamoClient := dynamo.InitClient(ctx)
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig.SavingGoalsTable)
		if err != nil {
			return
		}
	})

	return err
}

func (request *deleteSavingGoalRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeleteSavingGoalHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if dsgRequest == nil {
		dsgRequest = new(deleteSavingGoalRequest)
	}

	err := dsgRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_delete_saving_goal_request_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer dsgRequest.finish()

	return dsgRequest.process(ctx, req)
}

func (request *deleteSavingGoalRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingGoalID := req.PathParameters["savingGoalID"]

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	deleteSavingGoal := usecases.NewSavingGoalEliminator(request.savingGoalRepo)

	err = deleteSavingGoal(ctx, username, savingGoalID)
	if err != nil {
		logger.Error("delete_saving_goal_failed", err, req, models.Any("username", username))

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
