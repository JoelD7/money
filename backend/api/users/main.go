package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

type userRequest struct {
	*events.APIGatewayProxyRequest
	err error
}

var (
	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

func handler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	userID := req.QueryStringParameters["user_id"]

	user, err := storage.GetPersonByEmail(ctx, userID)
	if err != nil {
		return serverError(err)
	}

	if user == nil {
		return clientError(http.StatusNotFound)
	}

	personJson, err := json.Marshal(user)
	if err != nil {
		return serverError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(personJson),
	}, nil
}

func serverError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

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
