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
	userSaving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), userSaving)
	if err != nil {
		request.log.Error("request_body_unmarshal_failed", err, []models.LoggerObject{req})

		return getErrorResponse(errRequestBodyParseFailure)
	}

	updateSaving := usecases.NewSavingUpdater(request.savingsRepo)

	err = updateSaving(ctx, userSaving)
	if err != nil {
		request.log.Error("update_saving_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
	}, nil
}
