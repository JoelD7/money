package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var dsRequest *deleteSavingRequest
var dsOnce sync.Once

type deleteSavingRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
}

func (request *deleteSavingRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error
	dsOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig.SavingsTable, envConfig.PeriodSavingIndexName, envConfig.SavingGoalSavingIndexName)
		if err != nil {
			return
		}
		request.log = log
		request.log.SetHandler("delete-saving")
	})
	request.startingTime = time.Now()

	return err
}

func (request *deleteSavingRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeleteSavingHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if dsRequest == nil {
		dsRequest = new(deleteSavingRequest)
	}

	err := dsRequest.init(ctx, log, envConfig)
	if err != nil {
		log.Error("delete_saving_init_failed", err, []models.LoggerObject{req})
		return req.NewErrorResponse(err), nil
	}
	defer dsRequest.finish()

	return dsRequest.process(ctx, req)
}

func (request *deleteSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	userSaving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), userSaving)
	if err != nil {
		request.log.Error("request_body_unmarshal_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(models.ErrInvalidRequestBody), nil
	}

	deleteSaving := usecases.NewSavingDeleter(request.savingsRepo)

	err = deleteSaving(ctx, userSaving.SavingID, userSaving.Username)
	if err != nil {
		request.log.Error("delete_saving_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
