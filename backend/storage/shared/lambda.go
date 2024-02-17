package shared

import (
	"bytes"
	"context"
	"regexp"
	"runtime"
	"time"

	"github.com/JoelD7/money/backend/shared/env"
)

var (
	// LAMBDA_TIMEOUT is the environment variable that indicates the timeout for the lambda function.
	// https://docs.aws.amazon.com/lambda/latest/dg/configuration-function-common.html#configuration-timeout-console
	lambdaTimeout = env.GetString("LAMBDA_TIMEOUT", "10s")

	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)
)

// ExecuteLambda executes a lambda function's handler and returns the stack trace and error in case of a timeout.
func ExecuteLambda(handler func(ctx context.Context)) (string, error) {
	doneChan := make(chan struct{})

	ctx, cancel := getContextWithLambdaTimeout(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		handler(ctx)
		doneChan <- struct{}{}
	}(ctx)

	select {
	case <-ctx.Done():
		err := ctx.Err()

		if err != nil {
			clean := stackCleaner.FindAll(getStackTrace(), -1)

			return string(bytes.Join(clean, []byte("\n\n"))), err
		}

	case <-doneChan:
		return "", nil
	}

	return "", nil
}

// getContextWithLambdaTimeout returns a context with a timeout based on the invoking lambda function's timeout and a
// cancel function to cancel the context(https://pkg.go.dev/context#WithTimeout).
func getContextWithLambdaTimeout(parentCtx context.Context) (context.Context, context.CancelFunc) {
	duration, err := time.ParseDuration(lambdaTimeout)
	if err != nil {
		duration = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(parentCtx, duration)

	return ctx, cancel
}

// getStackTrace returns a formatted stack trace of all the goroutines executing in the program.
// This function is a variation of debug.Stack, but it returns the stack trace of all the goroutines. This behaviour
// is required in this case because the lambda execution occurs in a goroutine different from the one monitoring the
// context timeout.
func getStackTrace() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
