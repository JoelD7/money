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
			r.Get("/", getUserHandler)
			r.Route("/categories", func(r *router.Router) {
				r.Get("/", getCategoriesHandler)
				r.Post("/", createCategoryHandler)
				r.Put("/{categoryID}", updateCategoryHandler)
			})
		})

		r.Route("/savings", func(r *router.Router) {
			r.Get("/{savingID}", getSavingHandler)
			r.Get("/", getSavingsHandler)
			r.Post("/", createSavingHandler)
			r.Put("/{savingID}", updateSavingHandler)
			r.Delete("/", deleteSavingHandler)
		})

		r.Route("/periods", func(r *router.Router) {
			r.Post("/", createPeriodHandler)
			r.Put("/{periodID}", updatePeriodHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
