package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-generator/handler"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context) error {
		defer func() {
			err := logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return handler.Handle(ctx)
	})
}
