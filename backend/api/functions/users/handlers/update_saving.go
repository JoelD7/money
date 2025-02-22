package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var usRequest *updateSavingRequest
var usOnce sync.Once

type updateSavingRequest struct {
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
	periodRepo     period.Repository
}

func (request *updateSavingRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	usOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}

		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		logger.SetHandler("update-saving")
	})
	request.startingTime = time.Now()

	return err
}

func (request *updateSavingRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func UpdateSavingHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if usRequest == nil {
		usRequest = new(updateSavingRequest)
	}

	err := usRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("update_saving_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer usRequest.finish()

	return usRequest.process(ctx, req)
}

func (request *updateSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	userSaving, err := request.validateUpdateInputs(req)
	if err != nil {
		logger.Error("update_input_validation_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	updateSaving := usecases.NewSavingUpdater(request.savingsRepo, request.periodRepo, request.savingGoalRepo)

	saving, err := updateSaving(ctx, userSaving.Username, userSaving)
	if err != nil {
		logger.Error("update_saving_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, saving), nil
}

func (request *updateSavingRequest) validateUpdateInputs(req *apigateway.Request) (*models.Saving, error) {
	savingID, ok := req.PathParameters["savingID"]
	if !ok || savingID == "" {
		return nil, models.ErrMissingSavingID
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		return nil, models.ErrNoUsernameInContext
	}

	saving := &models.Saving{
		SavingID: savingID,
		Username: username,
	}

	err = json.Unmarshal([]byte(req.Body), saving)
	if err != nil {
		return nil, models.ErrInvalidRequestBody
	}

	err = validate.Email(username)
	if err != nil {
		return nil, err
	}

	err = validate.Amount(saving.Amount)
	if err != nil {
		return nil, err
	}

	return saving, nil
}
