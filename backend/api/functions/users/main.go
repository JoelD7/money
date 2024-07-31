package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/users/handlers"
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
		r.Route("/users", func(r *router.Router) {
			r.Get("/", handlers.GetUserHandler)
			r.Route("/categories", func(r *router.Router) {
				r.Get("/", handlers.GetCategoriesHandler)
				r.Post("/", handlers.CreateCategoryHandler)
				r.Put("/{categoryID}", handlers.UpdateCategoryHandler)
			})
		})

		r.Route("/savings", func(r *router.Router) {
			r.Get("/{savingID}", handlers.GetSavingHandler)
			r.Get("/", handlers.GetSavingsHandler)
			r.Post("/", handlers.CreateSavingHandler)
			r.Put("/{savingID}", handlers.UpdateSavingHandler)
			r.Delete("/", handlers.DeleteSavingHandler)
		})

		r.Route("/periods", func(r *router.Router) {
			r.Post("/", handlers.CreatePeriodHandler)
			r.Put("/{periodID}", handlers.UpdatePeriodHandler)
			r.Get("/{periodID}", handlers.GetPeriodHandler)
			r.Get("/", handlers.GetPeriodsHandler)
			r.Delete("/{periodID}", handlers.DeletePeriodHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
