package shared

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"time"
)

// LAMBDA_TIMEOUT is the environment variable that indicates the timeout for the lambda function.
// https://docs.aws.amazon.com/lambda/latest/dg/configuration-function-common.html#configuration-timeout-console
var lambdaTimeout = env.GetString("LAMBDA_TIMEOUT", "10s")

// GetContextWithLambdaTimeout returns a context with a timeout based on the invoking lambda function's timeout and a
// cancel function to cancel the context(https://pkg.go.dev/context#WithTimeout).
func GetContextWithLambdaTimeout(parentCtx context.Context, errChan chan error) (context.Context, context.CancelFunc) {
	//duration, err := time.ParseDuration(lambdaTimeout)
	//if err != nil {
	//	duration = 10 * time.Second
	//}

	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)

	go checkCtxError(ctx, errChan)

	return ctx, cancel
}

func checkCtxError(ctx context.Context, errChan chan error) {
	//TODO: what happens if context is never cancelled?
	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			fmt.Printf("checkCtxError: finished\n")
			return
		}
	}
}
