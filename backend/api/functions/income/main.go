package main

import (
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
)

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
