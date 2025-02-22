package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/users/handlers"
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
			r.Delete("/{savingID}", handlers.DeleteSavingHandler)

			r.Route("/goals", func(r *router.Router) {
				r.Post("/", handlers.CreateSavingGoalHandler)
				r.Get("/{savingGoalID}", handlers.GetSavingGoalHandler)
				r.Get("/", handlers.GetSavingGoalsHandler)
				r.Put("/{savingGoalID}", handlers.UpdateSavingGoalsHandler)
				r.Delete("/{savingGoalID}", handlers.DeleteSavingGoalHandler)
			})
		})

		r.Route("/periods", func(r *router.Router) {
			r.Post("/", handlers.CreatePeriodHandler)
			r.Get("/", handlers.GetPeriodsHandler)

			r.Route("/{periodID}", func(r *router.Router) {
				r.Put("/", handlers.UpdatePeriodHandler)
				r.Get("/", handlers.GetPeriodHandler)
				r.Delete("/", handlers.DeletePeriodHandler)

				r.Get("/stats", handlers.GetPeriodStatHandler)
			})
		})
	})

	lambda.Start(func(ctx context.Context, request *apigateway.Request) (res *apigateway.Response, err error) {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate(request.RequestContext.ExtendedRequestID))

		defer func() {
			err = logger.Finish()
			if err != nil {
				logger.ErrPrintln("failed to finish logger", err)
			}
		}()

		return rootRouter.Handle(ctx, request)
	})
}
