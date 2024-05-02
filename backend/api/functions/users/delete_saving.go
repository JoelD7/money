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

func (request *deleteSavingRequest) init(log logger.LogAPI) {
	dsOnce.Do(func() {
		dynamoClient := initDynamoClient()

		request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
		request.log = log
		request.log.SetHandler("delete-saving")
	})
	request.startingTime = time.Now()
}

func (request *deleteSavingRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func deleteSavingHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if dsRequest == nil {
		dsRequest = new(deleteSavingRequest)
	}

	dsRequest.init(log)
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
