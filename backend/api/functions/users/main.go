package main

import (
	"context"
	"errors"
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

	ErrNotFound = apigateway.NewError("not found", http.StatusNotFound)

	errUserFetchingFailed = errors.New("user fetching failed")
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

	rootRouter.Route("/users", func(r *router.Router) {
		r.Get("/{user-id}", getUserHandler)
	})

	lambda.Start(rootRouter.Handle)
}
