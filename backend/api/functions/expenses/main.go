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
		r.Route("/expenses", func(r *router.Router) {
			r.Get("/{expenseID}", getExpenseHandler)
			r.Put("/{expenseID}", updateExpenseHandler)
			r.Delete("/{expenseID}", deleteExpenseHandler)
			r.Get("/", getExpensesHandler)
			r.Post("/", createExpenseHandler)

			r.Route("/period", func(r *router.Router) {
				r.Patch("/missing", patchRecurrentExpenseHandler)
			})

		})
	})

	lambda.Start(rootRouter.Handle)
}
