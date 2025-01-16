package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
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
	startingTime time.Time
	err          error
	periodRepo   period.Repository
	cacheManager cache.IncomePeriodCacheManager
}

func (request *deletePeriodRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	dpOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)
		logger.SetHandler("delete-period")

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}

		request.cacheManager = cache.NewRedisCache()
	})
	request.startingTime = time.Now()

	return err
}

func (request *deletePeriodRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func DeletePeriodHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if dpRequest == nil {
		dpRequest = new(deletePeriodRequest)
	}

	err := dpRequest.init(ctx, envConfig)
	if err != nil {
		dpRequest.err = err

		logger.Error("delete_period_init_failed", err, req)

		return req.NewErrorResponse(err), nil

	}
	defer dpRequest.finish()

	return dpRequest.process(ctx, req)
}

func (request *deletePeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		logger.Error("missing_period_id", nil, req)

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	deletePeriod := usecases.NewPeriodDeleter(request.periodRepo, request.cacheManager)

	err = deletePeriod(ctx, periodID, username)
	if err != nil {
		logger.Error("delete_period_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusNoContent,
	}, nil
}
