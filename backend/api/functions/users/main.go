package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"time"
)

var (
	ErrNotFound = apigateway.NewError("not found", http.StatusNotFound)

	errUserFetchingFailed = errors.New("user fetching failed")
)

type userRequest struct {
	log          logger.LogAPI
	startingTime time.Time
	err          error
}

type userResponse struct {
	*models.User
	Remainder float64 `json:"remainder,omitempty"`
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

func handler(req *apigateway.Request) (*apigateway.Response, error) {
	request := &userRequest{
		log: logger.NewLogger(),
	}

	request.init()
	defer request.finish()

	return request.process(req)
}

func (request *userRequest) process(req *apigateway.Request) (*apigateway.Response, error) {
	ctx := context.Background()

	userID := req.PathParameters["user-id"]

	userRes, err := getUserResponse(ctx, userID)
	if err != nil {
		request.err = err
		request.log.Error("user_fetching_failed", err, []logger.Object{})

		return apigateway.NewErrorResponse(errUserFetchingFailed), nil
	}

	if userRes == nil {
		request.err = err
		request.log.Error("user_not_found", err, []logger.Object{})

		return apigateway.NewErrorResponse(ErrNotFound), nil
	}

	return apigateway.NewJSONResponse(http.StatusOK, userRes), nil
}

func getUserResponse(ctx context.Context, userID string) (*userResponse, error) {
	user, err := users.GetUserByEmail(ctx, userID)
	if err != nil {
		return nil, err
	}

	remainder, err := getUserPeriodRemainder(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userResponse{user, remainder}, nil
}

func getUserPeriodRemainder(ctx context.Context, user *models.User) (float64, error) {
	if user.CurrentPeriod == "" {
		return -1, nil
	}

	userExpenses, err := expenses.GetExpensesByPeriod(ctx, user.UserID, user.CurrentPeriod)
	if err != nil {
		return -1, err
	}

	userIncome, err := income.GetIncomeByPeriod(ctx, user.UserID, user.CurrentPeriod)
	if err != nil {
		return -1, err
	}

	totalExpense := 0.0

	for _, expense := range userExpenses {
		totalExpense += expense.Amount
	}

	totalIncome := 0.0
	for _, inc := range userIncome {
		totalIncome += inc.Amount
	}

	return totalIncome - totalExpense, nil
}

func main() {
	rootRouter := router.NewRouter()

	rootRouter.Route("/users", func(r *router.Router) {
		r.Get("/{user-id}", handler)
	})

	lambda.Start(rootRouter.Handle)
}
