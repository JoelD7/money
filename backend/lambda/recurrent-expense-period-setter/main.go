package main

import (
	"context"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-period-setter/handler"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	err := env.LoadEnv(context.Background())
	if err != nil {
		panic(err)
	}

	lambda.Start(handler.Handle)
}
