package handlers

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
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

func TestHandlerSuccess(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()

	usersMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	incomeMock := income.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)

	request := &getUserRequest{
		userRepo:     usersMock,
		expensesRepo: expensesMock,
		incomeRepo:   incomeMock,
		log:          logMock,
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}

	t.Run("Happy path", func(t *testing.T) {
		response, err := request.process(ctx, apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusOK, response.StatusCode)
		c.Contains(response.Body, `"remainder":8670`)
	})

	t.Run("No remainder set", func(t *testing.T) {
		mockedUser := users.GetDummyUser()
		mockedUser.CurrentPeriod = ""

		err := usersMock.CreateUser(ctx, mockedUser)
		c.NoError(err)

		response, err := request.process(ctx, apigwRequest)
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
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	request := &getUserRequest{
		userRepo:     usersMock,
		expensesRepo: expensesMock,
		incomeRepo:   incomeMock,
		log:          logMock,
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}

	t.Run("User fetching failed", func(t *testing.T) {
		dummyError := errors.New("get user failed")

		usersMock.ActivateForceFailure(dummyError)
		defer usersMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.NotNil(response)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		logMock.Output.Reset()
	})

	t.Run("User not found", func(t *testing.T) {
		usersMock.ActivateForceFailure(models.ErrUserNotFound)
		defer usersMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_fetching_failed")
		c.Contains(logMock.Output.String(), models.ErrUserNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Income not found", func(t *testing.T) {
		incomeMock.ActivateForceFailure(models.ErrIncomeNotFound)
		defer incomeMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusOK, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_remainder_could_not_be_calculated")
		c.Contains(logMock.Output.String(), models.ErrIncomeNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Expense not found", func(t *testing.T) {
		expensesMock.ActivateForceFailure(models.ErrExpensesNotFound)
		defer expensesMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.Nil(err)
		c.NotNil(response)
		c.Equal(http.StatusOK, response.StatusCode)
		c.Contains(logMock.Output.String(), "user_remainder_could_not_be_calculated")
		c.Contains(logMock.Output.String(), models.ErrExpensesNotFound.Error())
		logMock.Output.Reset()
	})
}
