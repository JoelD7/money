package shared

import (
	"bytes"
	"context"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/JoelD7/money/backend/shared/env"
)

var (
	stackCleaner = regexp.MustCompile(`[^\t]*:\d+`)
	errorLogger  = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

// ExecuteLambda executes a lambda function's code and returns the stack trace and error in case of a timeout.
// Takes the context(parentCtx) from inside lambda.Start and creates a new context with the timeout of the lambda.
// The handler function represents the lambda's code.
// For example: you would put something like this in the handler of the lambda:
//
//	  stackTrace, ctxError := shared.ExecuteLambda(func(ctx context.Context) {
//			res, err = req.process(ctx, event) this runs the lambda's code
//		})
//
//		if ctxError != nil {
//			req.log.Error("request_timeout", ctxError, []models.LoggerObject{
//				req.getEventAsLoggerObject(event),
//				req.log.MapToLoggerObject("stack", map[string]interface{}{
//					"s_trace": stackTrace,
//				}),
//			})
//		}
func ExecuteLambda(parentCtx context.Context, handler func(ctx context.Context)) (string, error) {
	doneChan := make(chan struct{})

	ctx, cancel := getContextWithLambdaTimeout(parentCtx)
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
	// LAMBDA_TIMEOUT is the environment variable that indicates the timeout for the lambda function in seconds.
	// https://docs.aws.amazon.com/lambda/latest/dg/configuration-function-common.html#configuration-timeout-console
	lambdaTimeout := env.GetString("LAMBDA_TIMEOUT", "")

	durationInt, err := strconv.Atoi(lambdaTimeout)
	if err != nil {
		errorLogger.Println("Error converting LAMBDA_TIMEOUT to int: ", err)
		return context.WithTimeout(parentCtx, 10*time.Second)
	}

	//Context timeout should be 1 second less than the lambda timeout to avoid the lambda being killed by AWS.
	durationInt--

	return context.WithTimeout(parentCtx, time.Duration(durationInt)*time.Second)
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
