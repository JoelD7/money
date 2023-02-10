// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/shared/router"
	"github.com/JoelD7/money/api/storage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	errMissingEmail     = errors.New("missing email")
	errMissingPassword  = errors.New("missing password")
	errWrongCredentials = errors.New("the email or password are incorrect")

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

const (
	passwordCost = bcrypt.DefaultCost
)

type signUpBody struct {
	FullName string `json:"fullname"`
	*Credentials
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signUpHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBody := &signUpBody{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateCredentials(reqBody.Credentials)
	if err != nil {
		return clientError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), passwordCost)
	if err != nil {
		return serverError(err)
	}

	fmt.Println("password: ", string(hashedPassword))

	err = storage.CreatePerson(reqBody.FullName, reqBody.Email, string(hashedPassword))
	if err != nil && errors.Is(err, storage.ErrExistingUser) {
		return clientError(err)
	}

	if err != nil {
		return serverError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Sign up succeeded",
	}, nil
}

func logInHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBody := &Credentials{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateCredentials(reqBody)
	if err != nil {
		return clientError(err)
	}

	person, err := storage.GetPerson(reqBody.Email)
	if err != nil {
		return clientError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(reqBody.Password))
	if err != nil {
		return clientError(errWrongCredentials)
	}

	//fmt.Println("FullName: ", reqBody.FullName, "Password hash: ", string(hashedPassword))

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Logged in!",
	}, nil
}

func serverError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	responseBody := strings.ToUpper(err.Error()[0:1]) + err.Error()[1:]

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       responseBody,
	}, nil
}

func validateCredentials(login *Credentials) error {
	if login.Email == "" {
		return errMissingEmail
	}

	if login.Password == "" {
		return errMissingPassword
	}

	return nil
}

func main() {
	route := router.NewRouter()

	route.Route("/auth", func(r *router.Router) {
		r.Post("/login", logInHandler)
		r.Post("/signup", signUpHandler)
	})

	lambda.Start(route.Handle)
}
