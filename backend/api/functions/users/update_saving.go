package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

type updateSavingRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
}

func (request *updateSavingRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *updateSavingRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func updateSavingHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(updateSavingRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *updateSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingID, ok := req.PathParameters["savingID"]
	if !ok {
		request.log.Error("missing_saving_id", errMissingSavingID, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(errMissingSavingID), nil
	}

	userSaving, err := validateUpdateInputs(req)
	if err != nil {
		request.log.Error("update_input_validation_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	userSaving.SavingID = savingID

	updateSaving := usecases.NewSavingUpdater(request.savingsRepo)

	err = updateSaving(ctx, userSaving)
	if err != nil {
		request.log.Error("update_saving_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func validateUpdateInputs(req *apigateway.Request) (*models.Saving, error) {
	saving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), saving)
	if err != nil {
		return nil, errRequestBodyParseFailure
	}

	err = validateEmail(saving.Username)
	if err != nil {
		return nil, err
	}

	err = validateAmount(saving.Amount)
	if err != nil {
		return nil, err
	}

	return saving, nil
}
