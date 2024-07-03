package main

import (
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/", func(r *router.Router) {
		r.Route("/expenses", func(r *router.Router) {
			r.Get("/{expenseID}", handlers.GetExpense)
			r.Put("/{expenseID}", handlers.UpdateExpense)
			r.Delete("/{expenseID}", handlers.DeleteExpense)
			r.Get("/", handlers.GetExpenses)
			r.Post("/", handlers.CreateExpense)

			r.Route("/period", func(r *router.Router) {
				r.Patch("/missing", handlers.PatchRecurrentExpense)
			})

		})
	})

	lambda.Start(rootRouter.Handle)
}
