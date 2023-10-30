package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"strconv"
)

var (
	awsRegion = env.GetString("REGION", "us-east-1")
)

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func getRequestQueryParams(req *apigateway.Request) (string, int, error) {
	pageSizeParam := 0
	var err error

	if req.QueryStringParameters["page_size"] != "" {
		pageSizeParam, err = strconv.Atoi(req.QueryStringParameters["page_size"])
		if err != nil {
			return "", 0, fmt.Errorf("%w: %v", models.ErrInvalidPageSize, err)
		}
	}

	return req.QueryStringParameters["start_key"], pageSizeParam, nil
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/", func(r *router.Router) {
		r.Route("/users", func(r *router.Router) {
			r.Get("/", getUserHandler)
			r.Route("/categories", func(r *router.Router) {
				r.Get("/", getCategoriesHandler)
				r.Post("/", createCategoryHandler)
				r.Put("/{categoryID}", updateCategoryHandler)
			})
		})

		r.Route("/savings", func(r *router.Router) {
			r.Get("/{savingID}", getSavingHandler)
			r.Get("/", getSavingsHandler)
			r.Post("/", createSavingHandler)
			r.Put("/{savingID}", updateSavingHandler)
			r.Delete("/", deleteSavingHandler)
		})

		r.Route("/periods", func(r *router.Router) {
			r.Post("/", createPeriodHandler)
			r.Put("/{periodID}", updatePeriodHandler)
			r.Get("/{periodID}", getPeriodHandler)
			r.Get("/", getPeriodsHandler)
		})
	})

	lambda.Start(rootRouter.Handle)
}
