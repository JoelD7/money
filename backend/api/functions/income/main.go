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
	awsRegion = env.GetString("AWS_REGION", "")
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
		r.Route("/income", func(r *router.Router) {
			r.Post("/", createIncomeHandler)
			r.Get("/{incomeID}", getIncomeHandler)
			r.Get("/", getMultipleIncomeHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
