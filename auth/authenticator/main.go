// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/shared/router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

var (
	errMissingUsername = errors.New("missing Username")
	errMissingPassword = errors.New("missing Password")

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

type loginBody struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func loginHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBody := &loginBody{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateParams(reqBody)
	if err != nil {
		return clientError(err)
	}

	fmt.Println("Username: ", reqBody.Username, "Password: ", reqBody.Password)

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Everything ok!",
	}, nil
}

func serverError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

func clientError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func validateParams(login *loginBody) error {
	if login.Username == "" {
		return errMissingUsername
	}

	if login.Password == "" {
		return errMissingPassword
	}

	return nil
}

func main() {
	route := router.NewRouter()

	route.Route("/auth", func(r *router.Router) {
		r.Post("/login", loginHandler)

	})

	lambda.Start(route.Handle)
}
