package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
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
		r.Route("/expenses", func(r *router.Router) {
			r.Get("/{expenseID}", handlers.GetExpense)
			r.Put("/{expenseID}", handlers.UpdateExpense)
			r.Delete("/{expenseID}", handlers.DeleteExpense)
			r.Get("/", handlers.GetExpenses)
			r.Post("/", handlers.CreateExpense)

			r.Route("/recurring", func(r *router.Router) {
				r.Delete("/{expenseRecurringID}", handlers.DeleteExpenseRecurring)
			})
		})
	})

	lambda.Start(rootRouter.Handle)
}
