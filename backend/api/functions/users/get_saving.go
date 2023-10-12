package main

import (
	"context"
	"errors"
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

var (
	errMissingSavingID = apigateway.NewError("missing saving ID", http.StatusBadRequest)
)

type getSavingRequest struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
}

func (request *getSavingRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.savingGoalRepo = savingoal.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getSavingRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getSavingHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getSavingRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *getSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingID, ok := req.PathParameters["savingID"]
	if !ok {
		request.log.Error("missing_saving_id", errMissingSavingID, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(errMissingSavingID), nil
	}

	username, err := getUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err, []models.LoggerObject{
			request.log.MapToLoggerObject("user_data", map[string]interface{}{
				"s_username": username,
			}),
		})

		return apigateway.NewErrorResponse(err), nil
	}

	getSaving := usecases.NewSavingGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	saving, err := getSaving(ctx, username, savingID)
	if errors.Is(err, models.ErrSavingGoalNameSettingFailed) {
		request.log.Error("get_saving_goal_name_failed", err, []models.LoggerObject{req})

		return apigateway.NewJSONResponse(http.StatusOK, saving), nil
	}

	if err != nil {
		request.log.Error("get_saving_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, saving), nil
}
