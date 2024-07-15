package main

import (
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-generator/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handle)
}
