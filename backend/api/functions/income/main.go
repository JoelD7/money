package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	envConfig, err := env.LoadEnv(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to load environment variables: %w", err))
	}

	rootRouter := router.NewRouter(envConfig)

	rootRouter.Route("/", func(r *router.Router) {
		r.Route("/income", func(r *router.Router) {
			r.Post("/", createIncomeHandler)
			r.Get("/{incomeID}", getIncomeHandler)
			r.Get("/", getMultipleIncomeHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
