package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/income/handlers"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/shared/uuid"
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
			r.Post("/", handlers.CreateIncomeHandler)
			r.Get("/{incomeID}", handlers.GetIncomeHandler)
			r.Get("/", handlers.GetMultipleIncomeHandler)
		})
	})

	lambda.Start(func(ctx context.Context, request *apigateway.Request) (res *apigateway.Response, err error) {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate(request.RequestContext.ExtendedRequestID))

		defer func() {
			err = logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return rootRouter.Handle(ctx, request)
	})
}
