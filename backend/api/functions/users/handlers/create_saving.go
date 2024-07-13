package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	csRequest *createSavingRequest
	csOnce    sync.Once

	periodTableNameEnv       = env.GetString("PERIOD_TABLE_NAME", "")
	uniquePeriodTableNameEnv = env.GetString("UNIQUE_PERIOD_TABLE_NAME", "")
	tableName                = env.GetString("SAVINGS_TABLE_NAME", "")
	periodSavingIndex        = env.GetString("PERIOD_SAVING_INDEX_NAME", "")
	savingGoalSavingIndex    = env.GetString("SAVING_GOAL_SAVING_INDEX_NAME", "")
)

type createSavingRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
	userRepo     users.Repository
	periodRepo   period.Repository
}

func (request *createSavingRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	csOnce.Do(func() {
		request.log = log
		request.log.SetHandler("create-saving")
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, tableName, periodSavingIndex, savingGoalSavingIndex)
		if err != nil {
			return
		}
		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *createSavingRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateSavingHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if csRequest == nil {
		csRequest = new(createSavingRequest)
	}

	err := csRequest.init(ctx, log)
	if err != nil {
		log.Error("init_create_saving_request_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}
	defer csRequest.finish()

	return csRequest.process(ctx, req)
}

func (request *createSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	userSaving, err := validateBody(req)
	if err != nil {
		request.log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createSaving := usecases.NewSavingCreator(request.savingsRepo, request.periodRepo)

	saving, err := createSaving(ctx, username, userSaving)
	if err != nil {
		request.log.Error("create_saving_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, saving), nil
}

func validateBody(req *apigateway.Request) (*models.Saving, error) {
	userSaving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), userSaving)
	if err != nil {
		return nil, models.ErrInvalidRequestBody
	}

	if userSaving.Amount == nil || (userSaving.Amount != nil && *userSaving.Amount == 0) {
		return nil, models.ErrMissingAmount
	}

	err = validate.Amount(userSaving.Amount)
	if err != nil {
		return nil, models.ErrInvalidSavingAmount
	}

	if userSaving.Period == nil || (userSaving.Period != nil && *userSaving.Period == "") {
		return nil, models.ErrMissingPeriod
	}

	return userSaving, nil
}
