package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

type userRequest struct {
	log *logger.Logger
}

func handler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	request := &userRequest{
		log: logger.NewLogger(),
	}

	return request.process(req)
}

func (request *userRequest) process(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	userID := req.QueryStringParameters["user_id"]

	user, err := storage.GetPersonByEmail(ctx, userID)
	if err != nil {
		request.log.Error("user_fetching_failed", err, []logger.Object{})

		return serverError()
	}

	if user == nil {
		request.log.Error("user_not_found", err, []logger.Object{})

		return clientError(http.StatusNotFound)
	}

	personJson, err := json.Marshal(user)
	if err != nil {
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
		r.Get("/", handler)
	})

	lambda.Start(rootRouter.Handle)
}
