package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage/person"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"time"
)

type userRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
}

func (request *userRequest) init() {
	request.startingTime = time.Now()
}

func (request *userRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func handler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	request := &userRequest{
		log: logger.NewLogger(),
	}

	request.init()
	defer request.finish()

	return request.process(req)
}

func (request *userRequest) process(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	userID := req.PathParameters["user-id"]

	user, err := person.GetPersonByEmail(ctx, userID)
	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, []logger.Object{})

		return serverError()
	}

	if user == nil {
		request.err = err
		request.log.Error("user_not_found", err, []logger.Object{})

		return clientError(http.StatusNotFound)
	}

	personJson, err := json.Marshal(user)
	if err != nil {
		request.err = err
		request.log.Error("user_response_marshal_failed", err, []logger.Object{})

		return serverError()
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(personJson),
	}, nil
}

func serverError() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/users", func(r *router.Router) {
		r.Get("/{user-id}", handler)
	})

	lambda.Start(rootRouter.Handle)
}
