package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	awsRegion = env.GetString("REGION", "us-east-1")
)

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
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
