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

var gpRequest *getPeriodRequest
var gpOnce sync.Once

type getPeriodRequest struct {
	startingTime time.Time
	err          error
	periodRepo   period.Repository
}

func (request *getPeriodRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	gpOnce.Do(func() {
		logger.SetHandler("get-period")
		dynamoClient := dynamo.InitClient(ctx)

		request.periodRepo, err = period.NewDynamoRepository(dynamoClient, envConfig.PeriodTable, envConfig.UniquePeriodTable)
		if err != nil {
			return
		}
	})
	request.startingTime = time.Now()

	return err
}

func (request *getPeriodRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetPeriodHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gpRequest == nil {
		gpRequest = new(getPeriodRequest)
	}

	err := gpRequest.init(ctx, envConfig)
	if err != nil {
		gpRequest.err = err

		logger.Error("get_period_init_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer gpRequest.finish()

	return gpRequest.process(ctx, req)
}

func (request *getPeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.err = models.ErrMissingPeriodID
		logger.Error("missing_period_id", request.err, req)

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		logger.Error("get_username_from_context_failed", request.err, req)

		return req.NewErrorResponse(err), nil
	}

	getPeriod := usecases.NewPeriodGetter(request.periodRepo)

	userPeriod, err := getPeriod(ctx, username, periodID)
	if err != nil {
		request.err = err
		logger.Error("get_period_failed", request.err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, userPeriod), nil
}
