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
		models.ErrUserNotFound:         {HTTPCode: http.StatusNotFound, Message: models.ErrUserNotFound.Error()},
		models.ErrIncomeNotFound:       {HTTPCode: http.StatusNotFound, Message: models.ErrIncomeNotFound.Error()},
		models.ErrExpensesNotFound:     {HTTPCode: http.StatusNotFound, Message: models.ErrExpensesNotFound.Error()},
		errNoUserEmailInContext:        {HTTPCode: http.StatusBadRequest, Message: errNoUserEmailInContext.Error()},
		errRequestBodyParseFailure:     {HTTPCode: http.StatusBadRequest, Message: errRequestBodyParseFailure.Error()},
		models.ErrSavingsNotFound:      {HTTPCode: http.StatusNotFound, Message: models.ErrSavingsNotFound.Error()},
		models.ErrInvalidAmount:        {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidAmount.Error()},
		models.ErrMissingUsername:      {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingUsername.Error()},
		models.ErrInvalidEmail:         {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidEmail.Error()},
		models.ErrInvalidRequestBody:   {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidRequestBody.Error()},
		models.ErrMissingSavingID:      {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingSavingID.Error()},
		models.ErrUpdateSavingNotFound: {HTTPCode: http.StatusNotFound, Message: models.ErrUpdateSavingNotFound.Error()},
		models.ErrDeleteSavingNotFound: {HTTPCode: http.StatusNotFound, Message: models.ErrDeleteSavingNotFound.Error()},
		models.ErrInvalidPageSize:      {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidPageSize.Error()},
		models.ErrInvalidStartKey:      {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidStartKey.Error()},
		errMissingSavingID:             {HTTPCode: http.StatusBadRequest, Message: errMissingSavingID.Error()},
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

	return apigateway.NewErrorResponse(err), nil
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/", func(r *router.Router) {
		r.Route("/users", func(r *router.Router) {
			r.Get("/{username}", getUserHandler)
		})

		r.Route("/savings", func(r *router.Router) {
			r.Get("/{savingID}", getSavingHandler)
			r.Get("/", getSavingsHandler)
			r.Post("/", createSavingHandler)
			r.Put("/", updateSavingHandler)
			r.Delete("/", deleteSavingHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
