package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gpRequest *getPeriodRequest
var gpOnce sync.Once

type getPeriodRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	periodRepo   period.Repository
	usersRepo    users.Repository
}

func (request *getPeriodRequest) init(ctx context.Context, log logger.LogAPI) {
	gpOnce.Do(func() {
		dynamoClient := dynamo.InitDynamoClient(ctx)

		request.periodRepo = period.NewDynamoRepository(dynamoClient)
		request.usersRepo = users.NewDynamoRepository(dynamoClient)
		request.log = log
		request.log.SetHandler("get-period")
	})
	request.startingTime = time.Now()
}

func (request *getPeriodRequest) finish() {
	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetPeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if gpRequest == nil {
		gpRequest = new(getPeriodRequest)
	}

	gpRequest.init(ctx, log)
	defer gpRequest.finish()

	return gpRequest.process(ctx, req)
}

func (request *getPeriodRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.err = models.ErrMissingPeriodID
		request.log.Error("missing_period_id", request.err, []models.LoggerObject{req})

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.log.Error("get_username_from_context_failed", request.err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	getPeriod := usecases.NewPeriodGetter(request.periodRepo, request.usersRepo)

	userPeriod, err := getPeriod(ctx, username, periodID)
	if err != nil {
		request.err = err
		request.log.Error("get_period_failed", request.err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, userPeriod), nil
}
