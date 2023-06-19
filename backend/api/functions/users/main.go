package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"time"
)

var (
	ErrNotFound = apigateway.NewError("not found", http.StatusNotFound)
)

type userRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
}

func (request *userRequest) init() {
	request.startingTime = time.Now()
}

func (request *userRequest) finish() {
	defer func() {
		err := request.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	request.log.LogLambdaTime(request.startingTime, request.err, recover())
}

func handler(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	request := &userRequest{
		log: logger.NewLogger(),
	}

	request.init()
	defer request.finish()

	return request.process(req)
}

func (request *userRequest) process(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	userID := req.PathParameters["user-id"]

	user, err := users.GetUserByEmail(ctx, userID)
	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, []logger.Object{})

		return apigateway.NewErrorResponse(nil), nil
	}

	if user == nil {
		request.err = err
		request.log.Error("user_not_found", err, []logger.Object{})

		return apigateway.NewErrorResponse(ErrNotFound), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, user), nil
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/users", func(r *router.Router) {
		r.Get("/{user-id}", handler)
	})

	lambda.Start(rootRouter.Handle)
}
