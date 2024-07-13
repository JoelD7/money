package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var dpRequest *deletePeriodRequest
var dpOnce sync.Once

type deletePeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *deletePeriodRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	dpOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)
		request.log = log
		request.log.SetHandler("delete-period")

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *deletePeriodRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeletePeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if dpRequest == nil {
		dpRequest = new(deletePeriodRequest)
	}

	err := dpRequest.init(ctx, log)
	if err != nil {
		dpRequest.err = err

		log.Error("delete_period_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil

	}
	defer dpRequest.finish()

	return dpRequest.process(ctx, req)
}

func (request *deletePeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.log.Error("missing_period_id", nil, []models.LoggerObject{req})

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	deletePeriod := usecases.NewPeriodDeleter(request.periodRepo)

	err = deletePeriod(ctx, periodID, username)
	if err != nil {
		request.log.Error("delete_period_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
