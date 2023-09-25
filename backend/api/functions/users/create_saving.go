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

var (
	errRequestBodyParseFailure = apigateway.NewError("couldn't parse the request body. Please check it", http.StatusBadRequest)
)

type createSavingRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
}

func (request *createSavingRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *createSavingRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func createSavingHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(createSavingRequest)

	request.init()
	defer request.finish()

	return request.process(ctx, req)
}

func (request *createSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	userSaving := new(models.Saving)

	err := json.Unmarshal([]byte(req.Body), userSaving)
	if err != nil {
		request.log.Error("request_body_unmarshal_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(errRequestBodyParseFailure), nil
	}

	createSaving := usecases.NewSavingCreator(request.savingsRepo, request.log)

	err = createSaving(ctx, userSaving)
	if err != nil {
		request.log.Error("create_saving_failed", err, []models.LoggerObject{req})

		return apigateway.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusCreated,
	}, nil
}
