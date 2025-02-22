package main

import (
	"context"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-period-setter/handler"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/uuid"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate("recurrent-expense-period-setter"))

		defer func() {
			err := logger.Finish()
			if err != nil {
				logger.ErrPrintln("failed to finish logger", err)
			}
		}()

		return handler.Handle(ctx, sqsEvent)
	})
}
