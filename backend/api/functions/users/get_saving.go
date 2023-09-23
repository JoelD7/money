package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"time"
)

var (
	errMissingSavingID = errors.New("missing savingID")
)

type getSavingRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
}

func (request *getSavingRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
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
		request.log.Error("saving_id", errMissingSavingID, []models.LoggerObject{req})

		return getErrorResponse(errMissingSavingID)
	}

	username, err := getUsernameFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	getSaving := usecases.NewSavingGetter(request.savingsRepo, request.log)

	saving, err := getSaving(ctx, username, savingID)
	if err != nil {
		request.log.Error("get_saving_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	responseJSON, err := json.Marshal(saving)
	if err != nil {
		request.log.Error("get_saving_marshal_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}
