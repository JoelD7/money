package main

import (
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/api/storage"
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

	userResources = []string{"savings", "categories", "month-budget"}
)

func apiGatewayHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := &userRequest{APIGatewayProxyRequest: &request}

	return router(req)
}

func router(req *userRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "/users/{user-id}":
		js, _ := json.Marshal(req)
		fmt.Println("req: ", string(js))
		return events.APIGatewayProxyResponse{}, nil
	case http.MethodGet
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "This method is not supported",
		}, nil
	}
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.QueryStringParameters["user_id"]

	person, err := storage.GetPerson(userID)
	if err != nil {
		return serverError(err)
	}

	if person == nil {
		return clientError(http.StatusNotFound)
	}

	personJson, err := json.Marshal(person)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(personJson),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(apiGatewayHandler)
}
