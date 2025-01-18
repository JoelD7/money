package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-period-setter/handler"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	_, err := env.LoadEnv(context.Background())
	if err != nil {
		panic(err)
	}

	logger.InitLogger(logger.LogstashImplementation)

	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		defer func() {
			err = logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return handler.Handle(ctx, sqsEvent)
	})
}
