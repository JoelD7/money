package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

var cpRequest *CreatePeriodRequest
var cpOnce sync.Once

type CreatePeriodRequest struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	PeriodRepo   period.Repository
}

func (request *CreatePeriodRequest) init(ctx context.Context, log logger.LogAPI) error {
	var err error
	cpOnce.Do(func() {
		dynamoClient := dynamo.InitClient(ctx)

		request.PeriodRepo, err = period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
		if err != nil {
			return
		}
		request.Log = log
		request.Log.SetHandler("create-period")
	})
	request.startingTime = time.Now()

	return err
}

func (request *CreatePeriodRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreatePeriodHandler(ctx context.Context, log logger.LogAPI, req *apigateway.Request) (*apigateway.Response, error) {
	if cpRequest == nil {
		cpRequest = new(CreatePeriodRequest)
	}

	err := cpRequest.init(ctx, log)
	if err != nil {
		log.Error("create_period_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}
	defer cpRequest.finish()

	return cpRequest.Process(ctx, req)
}

func (request *CreatePeriodRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodModel, err := request.validateCreateRequestBody(req)
	if err != nil {
		request.err = err
		request.Log.Error("validate_request_body_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.Log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	createPeriod := usecases.NewPeriodCreator(request.PeriodRepo, request.Log)

	createdPeriod, err := createPeriod(ctx, username, periodModel)
	if err != nil {
		request.err = err
		request.Log.Error("create_period_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, createdPeriod), nil
}

func (request *CreatePeriodRequest) validateCreateRequestBody(req *apigateway.Request) (*models.Period, error) {
	p := new(models.Period)

	err := json.Unmarshal([]byte(req.Body), p)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, models.ErrInvalidRequestBody)
	}

	if p.Name == nil || p.Name != nil && *p.Name == "" {
		return nil, models.ErrMissingPeriodName
	}

	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return nil, models.ErrMissingPeriodDates
	}

	return p, nil
}
