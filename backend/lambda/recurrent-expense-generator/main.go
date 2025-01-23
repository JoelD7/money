package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-generator/handler"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/uuid"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context) error {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate("recurrent-expense-generator"))

		defer func() {
			err := logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return handler.Handle(ctx)
	})
}
