package handlers

import (
	"context"
	"errors"
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

var (
	errMissingSavingID = apigateway.NewError("missing saving ID", http.StatusBadRequest)

	gsRequest *getSavingRequest
	gsOnce    sync.Once
)

type getSavingRequest struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
}

func (request *getSavingRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gsOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig.SavingsTable, envConfig.PeriodSavingIndexName, envConfig.SavingGoalSavingIndexName)
		if err != nil {
			return
		}
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig.SavingGoalsTable)
		if err != nil {
			return
		}
		request.log = log
		request.log.SetHandler("get-saving")
	})
	request.startingTime = time.Now()

	return err
}

func (request *getSavingRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetSavingHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gsRequest == nil {
		gsRequest = new(getSavingRequest)
	}

	err := gsRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("get_saving_init_failed", err, req)

		return req.NewErrorResponse(err), nil

	}
	defer gsRequest.finish()

	return gsRequest.process(ctx, req)
}

func (request *getSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingID, ok := req.PathParameters["savingID"]
	if !ok {
		request.log.Error("missing_saving_id", errMissingSavingID, req)

		return req.NewErrorResponse(errMissingSavingID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	err = validate.Email(username)
	if err != nil {
		request.log.Error("invalid_username", err,
			request.log.MapToLoggerObject("user_data", map[string]interface{}{
				"s_username": username,
			}),
		)

		return req.NewErrorResponse(err), nil
	}

	getSaving := usecases.NewSavingGetter(request.savingsRepo, request.savingGoalRepo, request.log)

	saving, err := getSaving(ctx, username, savingID)
	if errors.Is(err, models.ErrSavingGoalNameSettingFailed) {
		request.log.Error("get_saving_goal_name_failed", err, req)

		return req.NewJSONResponse(http.StatusOK, saving), nil
	}

	if err != nil {
		request.log.Error("get_saving_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, saving), nil
}
