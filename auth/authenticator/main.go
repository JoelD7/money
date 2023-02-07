// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/shared/router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var (
	errMissingUsername = errors.New("missing username")
	errMissingPassword = errors.New("missing password")
)

func loginHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	username, ok := request.QueryStringParameters["username"]
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "The username is required",
		}, errMissingUsername
	}

	password := request.QueryStringParameters["password"]
	if !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "The password is required",
		}, errMissingPassword
	}

	fmt.Println("username: ", username, "password: ", password)

	return &events.APIGatewayProxyResponse{}, nil
}

func main() {
	route := router.NewRouter()

	route.Route("/auth", func(r *router.Router) {
		r.Post("/login", loginHandler)

	})

	lambda.Start(route.Handle)
}
