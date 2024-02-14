package shared

import (
	"context"
	"github.com/JoelD7/money/backend/shared/env"
	"time"
)

// LAMBDA_TIMEOUT is the environment variable that indicates the timeout for the lambda function.
// https://docs.aws.amazon.com/lambda/latest/dg/configuration-function-common.html#configuration-timeout-console
var lambdaTimeout = env.GetString("LAMBDA_TIMEOUT", "10s")

// GetContextWithLambdaTimeout returns a context with a timeout based on the invoking lambda function's timeout.
func GetContextWithLambdaTimeout(parentCtx context.Context) (context.Context, context.CancelFunc) {
	duration, err := time.ParseDuration(lambdaTimeout)
	if err != nil {
		duration = 10 * time.Second
	}

	return context.WithTimeout(parentCtx, duration)
}
