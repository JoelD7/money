package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
)

var (
	awsRegion = env.GetString("REGION", "us-east-1")

	responseByErrors = map[error]apigateway.Error{
		models.ErrUserNotFound:     {HTTPCode: http.StatusNotFound, Message: models.ErrUserNotFound.Error()},
		models.ErrIncomeNotFound:   {HTTPCode: http.StatusNotFound, Message: models.ErrIncomeNotFound.Error()},
		models.ErrExpensesNotFound: {HTTPCode: http.StatusNotFound, Message: models.ErrExpensesNotFound.Error()},
		errNoUserEmailInContext:    {HTTPCode: http.StatusBadRequest, Message: errNoUserEmailInContext.Error()},
		errRequestBodyParseFailure: {HTTPCode: http.StatusBadRequest, Message: errRequestBodyParseFailure.Error()},
		models.ErrSavingsNotFound:  {HTTPCode: http.StatusNotFound, Message: models.ErrSavingsNotFound.Error()},
		models.ErrInvalidAmount:    {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidAmount.Error()},
		models.ErrMissingEmail:     {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingEmail.Error()},
		models.ErrInvalidEmail:     {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidEmail.Error()},
	}
)

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func getErrorResponse(err error) (*apigateway.Response, error) {
	for mappedErr, responseErr := range responseByErrors {
		if errors.Is(err, mappedErr) {
			return apigateway.NewJSONResponse(responseErr.HTTPCode, responseErr.Message), nil
		}
	}

	return apigateway.NewErrorResponse(err), err
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/", func(r *router.Router) {
		r.Route("/users", func(r *router.Router) {
			r.Get("/{user-id}", getUserHandler)
		})

		r.Route("/savings", func(r *router.Router) {
			r.Get("/", getSavingsHandler)
			r.Post("/", createSavingHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
