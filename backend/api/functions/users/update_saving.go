package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type updateSavingRequest struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	savingsRepo    savings.Repository
	savingGoalRepo savingoal.Repository
	periodRepo     period.Repository
}

func (request *updateSavingRequest) init(log logger.LogAPI) {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.periodRepo = period.NewDynamoRepository(dynamoClient)
	request.savingGoalRepo = savingoal.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = log
	request.log.SetHandler("update-saving")
}

func (request *updateSavingRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updateSavingHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updateSavingRequest)

	request.init(log)
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updateSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	userSaving, err := request.validateUpdateInputs(req)
	if err != nil {
		request.log.Error("update_input_validation_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	updateSaving := usecases.NewSavingUpdater(request.savingsRepo, request.periodRepo, request.savingGoalRepo)

	saving, err := updateSaving(ctx, userSaving.Username, userSaving)
	if err != nil {
		request.log.Error("update_saving_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, saving), nil
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
