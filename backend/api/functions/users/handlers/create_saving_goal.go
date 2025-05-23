package handlers

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var (
	csgRequest *createSavingGoalRequest
	csgOnce    sync.Once
)

type createSavingGoalRequest struct {
	startingTime   time.Time
	err            error
	savingGoalRepo savingoal.Repository
	cacheManager   cache.IdempotenceCacheManager
}

func (request *createSavingGoalRequest) init(ctx context.Context, envConfig *models.EnvironmentConfiguration) error {
	var err error
	csgOnce.Do(func() {
		request.startingTime = time.Now()
		dynamoClient := dynamo.InitClient(ctx)
		request.savingGoalRepo, err = savingoal.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
		request.cacheManager = cache.NewRedisCache()
	})

	return err
}

func (request *createSavingGoalRequest) finish() {
	logger.LogLambdaTime(request.startingTime, request.err, recover())
}

func CreateSavingGoalHandler(ctx context.Context, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if csgRequest == nil {
		csgRequest = new(createSavingGoalRequest)
	}

	err := csgRequest.init(ctx, envConfig)
	if err != nil {
		logger.Error("init_create_saving_goal_request_failed", err, req)

		return req.NewErrorResponse(err), nil
	}
	defer csgRequest.finish()

	return csgRequest.process(ctx, req)
}

func (request *createSavingGoalRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	idempotencyKey, err := req.GetIdempotenceyKeyFromHeader()
	if err != nil {
		request.err = err
		logger.Error("http_request_validation_failed", err, req)
		return req.NewErrorResponse(err), nil
	}

	savingGoal, err := validateSavingGoalBody(req)
	if err != nil {
		logger.Error("validate_request_body_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		logger.Error("get_username_from_context_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	createSavingGoal := usecases.NewSavingGoalCreator(request.savingGoalRepo, request.cacheManager)

	savingGoal, err = createSavingGoal(ctx, username, idempotencyKey, savingGoal)
	if err != nil {
		logger.Error("create_saving_goal_failed", err, req)

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusCreated, savingGoal), nil
}

func validateSavingGoalBody(req *apigateway.Request) (*models.SavingGoal, error) {
	var savingGoal models.SavingGoal

	err := json.Unmarshal([]byte(req.Body), &savingGoal)
	if err != nil {
		return nil, err
	}

	if savingGoal.Name == nil || (savingGoal.Name != nil && *savingGoal.Name == "") {
		return nil, models.ErrMissingSavingGoalName
	}

	if savingGoal.Target == nil {
		return nil, models.ErrMissingSavingGoalTarget
	}

	if savingGoal.Target != nil && *savingGoal.Target <= 0 {
		return nil, models.ErrInvalidSavingGoalTarget
	}

	if savingGoal.Deadline != nil && savingGoal.Deadline.Before(time.Now()) {
		return nil, models.ErrInvalidSavingGoalDeadline
	}

	if savingGoal.IsRecurring && savingGoal.RecurringAmount == nil {
		return nil, models.ErrMissingSavingGoalRecurringAmount
	}

	if err = validate.Amount(savingGoal.RecurringAmount); savingGoal.IsRecurring && err != nil {
		return nil, err
	}

	return &savingGoal, nil
}
