package handlers

import (
	"context"
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
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
}

func (request *deleteSavingRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	dsOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.savingsRepo, err = savings.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
		logger.SetHandler("delete-saving")
	})
	request.startingTime = time.Now()

	return err
}

func (request *deleteSavingRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeleteSavingHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if dsRequest == nil {
		dsRequest = new(deleteSavingRequest)
	}

	err := dsRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("delete_saving_init_failed", err, req)
		return req.NewErrorResponse(err), nil
	}
	defer dsRequest.finish()

	return dsRequest.process(ctx, req)
}

func (request *deleteSavingRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	savingID := req.PathParameters["savingID"]

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	deleteSaving := usecases.NewSavingDeleter(request.savingsRepo)

	err = deleteSaving(ctx, savingID, username)
	if err != nil {
		logger.Error("delete_saving_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusNoContent, nil), nil
}
