package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
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
		r.Route("/expenses", func(r *router.Router) {
			r.Get("/{expenseID}", handlers.GetExpense)
			r.Put("/{expenseID}", handlers.UpdateExpense)
			r.Delete("/{expenseID}", handlers.DeleteExpense)
			r.Get("/", handlers.GetExpenses)
			r.Post("/", handlers.CreateExpense)

			r.Route("/recurring", func(r *router.Router) {
				r.Delete("/{expenseRecurringID}", handlers.DeleteExpenseRecurring)
			})

			r.Route("/stats", func(r *router.Router) {
				r.Route("/period", func(r *router.Router) {
					r.Get("/{periodID}", handlers.GetExpensesStats)
				})
			})
		})
	})

	lambda.Start(func(ctx context.Context, request *apigateway.Request) (res *apigateway.Response, err error) {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", request.RequestContext.ExtendedRequestID)

		defer func() {
			err = logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return rootRouter.Handle(ctx, request)
	})
}
