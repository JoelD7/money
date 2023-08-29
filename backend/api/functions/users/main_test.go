package main

import (
	"errors"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var (
	logMock *logger.LogMock
)

func init() {
	logMock = logger.InitLoggerMock(nil)
}

func TestHandlerSuccess(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	incomeMock := income.NewDynamoMock()

	request := &userRequest{
		userRepo:     users.NewRepository(usersMock),
		expensesRepo: expenses.NewRepository(expensesMock),
		incomeRepo:   income.NewRepository(incomeMock),
		log:          logger.NewLogger(),
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{},
		PathParameters: map[string]string{
			"user-id": "test@gmail.com",
		},
	}

	t.Run("Happy path", func(t *testing.T) {
		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusOK, response.StatusCode)
		c.Contains(response.Body, `"remainder":8670`)
	})

	t.Run("No remainder set", func(t *testing.T) {
		mockedUser := users.GetDummyUser()
		mockedUser.CurrentPeriod = ""

		usersMock.SetMockedUser(mockedUser)

		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusOK, response.StatusCode)
		c.NotContains(response.Body, "remainder")
	})
}

func TestHandlerFailed(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	incomeMock := income.NewDynamoMock()

	request := &userRequest{
		userRepo:     users.NewRepository(usersMock),
		expensesRepo: expenses.NewRepository(expensesMock),
		incomeRepo:   income.NewRepository(incomeMock),
		log:          logger.NewLogger(),
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{},
		PathParameters: map[string]string{
			"user-id": "test@gmail.com",
		},
	}

	t.Run("User fetching failed", func(t *testing.T) {
		usersMock.ActivateForceFailure(errors.New("get user failed"))
		defer usersMock.DeactivateForceFailure()

		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		logMock.Output.Reset()
	})

	t.Run("User not found", func(t *testing.T) {
		usersMock.SetMockedUser(nil)
		defer usersMock.SetMockedUser(users.GetDummyUser())

		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		c.Contains(logMock.Output.String(), users.ErrNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Income not found", func(t *testing.T) {
		incomeMock.SetMockedIncome(nil)
		defer incomeMock.SetMockedIncome(income.GetDummyIncome())

		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		c.Contains(logMock.Output.String(), income.ErrNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Expense not found", func(t *testing.T) {
		expensesMock.SetMockedExpenses(nil)
		defer expensesMock.SetMockedExpenses(expenses.GetDummyExpenses())

		response, err := request.process(apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		c.Contains(logMock.Output.String(), expenses.ErrNotFound.Error())
		logMock.Output.Reset()
	})
}