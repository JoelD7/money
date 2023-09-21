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
	"strconv"
	"time"
)

var (
	errNoUserEmailInContext = errors.New("couldn't identify the user to get the savings from. Check if your Bearer token header is correct")
)

type getSavingsRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	savingsRepo  savings.Repository
	username     string
}

type savingsResponse struct {
	Savings []*models.Saving `json:"savings"`
	NextKey string           `json:"next_key"`
}

func (request *getSavingsRequest) init() {
	dynamoClient := initDynamoClient()

	request.savingsRepo = savings.NewDynamoRepository(dynamoClient)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func (request *getSavingsRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func getSavingsHandler(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	request := new(getSavingsRequest)

	request.init()
	defer request.finish()

	err := request.setUsername(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	return request.routeToHandlers(ctx, req)
}

func (request *getSavingsRequest) setUsername(req *apigateway.Request) error {
	username, err := getUsernameFromContext(req)
	if err != nil {
		return err
	}

	request.username = username

	return nil
}

func (request *getSavingsRequest) routeToHandlers(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	_, periodOk := req.QueryStringParameters["period"]
	_, savingGoalOk := req.QueryStringParameters["saving_goal_id"]

	if periodOk && savingGoalOk {
		return request.getUserSavingsByPeriodAndSavingGoal(ctx, req)
	}

	if periodOk && !savingGoalOk {
		return request.getUserSavingsByPeriod(ctx, req)
	}

	if !periodOk && savingGoalOk {
		return request.getUserSavingsByPeriodAndSavingGoal(ctx, req)
	}

	return request.getUserSavings(ctx, req)
}

func (request *getSavingsRequest) getUserSavings(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getSavings := usecases.NewSavingsGetter(request.savingsRepo, request.log)

	startKey, pageSize, err := getRequestParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	userSavings, nextKey, err := getSavings(ctx, request.username, startKey, pageSize)
	if err != nil {
		request.log.Error("savings_fetch_failed", err, []models.LoggerObject{
			req,
		})

		return getErrorResponse(err)
	}

	response := &savingsResponse{userSavings, nextKey}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		request.log.Error("savings_marshal_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}, nil
}

func (request *getSavingsRequest) getUserSavingsByPeriod(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
}

func (request *getSavingsRequest) getUserSavingsBySavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
}

func (request *getSavingsRequest) getUserSavingsByPeriodAndSavingGoal(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
}

func getUsernameFromContext(req *apigateway.Request) (string, error) {
	username, ok := req.RequestContext.Authorizer["username"].(string)
	if !ok {
		return "", errNoUserEmailInContext
	}

	return username, nil
}

func getRequestParams(req *apigateway.Request) (string, int, error) {
	pageSizeParam := 0
	var err error

	if req.QueryStringParameters["page_size"] != "" {
		pageSizeParam, err = strconv.Atoi(req.QueryStringParameters["page_size"])
		if err != nil {
			return "", 0, err
		}
	}

	return req.QueryStringParameters["start_key"], pageSizeParam, nil
}
