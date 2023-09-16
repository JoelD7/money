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

	return request.process(ctx, req)
}

func (request *getSavingsRequest) process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	getSavings := usecases.NewSavingsGetter(request.savingsRepo, request.log)

	email, err := getUserEmailFromContext(req)
	if err != nil {
		request.log.Error("get_user_email_from_context_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	startKey, pageSize, err := getRequestParams(req)
	if err != nil {
		request.log.Error("get_request_params_failed", err, []models.LoggerObject{req})

		return getErrorResponse(err)
	}

	userSavings, nextKey, err := getSavings(ctx, email, startKey, pageSize)
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

func getUserEmailFromContext(req *apigateway.Request) (string, error) {
	email, ok := req.RequestContext.Authorizer["email"].(string)
	if !ok {
		return "", errNoUserEmailInContext
	}

	return email, nil
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
