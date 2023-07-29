package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/usecases"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
	"time"
)

var (
	awsRegion = env.GetString("REGION", "us-east-1")

	ErrNotFound = apigateway.NewError("not found", http.StatusNotFound)

	errUserFetchingFailed = errors.New("user fetching failed")
)

type userRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
	userRepo     *users.Repository
	incomeRepo   *income.Repository
	expensesRepo *expenses.Repository
}

func (request *userRequest) init() {
	dynamoClient := initDynamoClient()

	dynamoUserRepository := users.NewDynamoRepository(dynamoClient)
	dynamoIncomeRepository := income.NewDynamoRepository(dynamoClient)
	dynamoExpensesRepository := expenses.NewDynamoRepository(dynamoClient)

	request.userRepo = users.NewRepository(dynamoUserRepository)
	request.incomeRepo = income.NewRepository(dynamoIncomeRepository)
	request.expensesRepo = expenses.NewRepository(dynamoExpensesRepository)
	request.startingTime = time.Now()
	request.log = logger.NewLogger()
}

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
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

func handler(req *apigateway.Request) (*apigateway.Response, error) {
	request := new(userRequest)

	request.init()
	defer request.finish()

	return request.process(req)
}

func (request *userRequest) process(req *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	userID := req.PathParameters["user-id"]

	getUser := usecases.NewUserGetter(request.userRepo, request.incomeRepo, request.expensesRepo)

	user, err := getUser(ctx, userID)
	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, []logger.Object{})

		return apigateway.NewErrorResponse(errUserFetchingFailed), nil
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
