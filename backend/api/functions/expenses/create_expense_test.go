package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"net/http"
	"testing"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestHandlerSuccess(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	userMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	ctx := context.Background()

	request := &createExpenseRequest{
		log:          logMock,
		userRepo:     userMock,
		expensesRepo: expensesMock,
	}

	apigwRequest := getCreateExpenseRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func TestHandlerFailure(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	userMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	ctx := context.Background()

	request := &createExpenseRequest{
		log:          logMock,
		userRepo:     userMock,
		expensesRepo: expensesMock,
	}

	apigwRequest := getCreateExpenseRequest()

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest.Body = "invalid"
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Invalid email", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer["username"] = "invalid"
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_input_failed")
		logMock.Output.Reset()
	})

	t.Run("Missing name", func(t *testing.T) {
		apigwRequest.Body = `{"amount":893,"period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, models.ErrMissingName.Error())
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_input_failed")
		logMock.Output.Reset()
	})

	t.Run("Missing amount", func(t *testing.T) {
		apigwRequest.Body = `{"name":"Jordan shopping","period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, models.ErrMissingAmount.Error())
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_input_failed")
		logMock.Output.Reset()
	})

	t.Run("Invalid amount", func(t *testing.T) {
		apigwRequest.Body = `{"amount":0,"name":"Jordan shopping","period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, models.ErrInvalidAmount.Error())
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_input_failed")
		logMock.Output.Reset()
	})

	t.Run("Create expense failed", func(t *testing.T) {
		expensesMock.ActivateForceFailure(errors.New("dummy"))
		defer expensesMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "create_expense_failed")
		logMock.Output.Reset()
	})

	t.Run("Get username from context failed", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = nil
		defer func() { apigwRequest = getCreateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_username_from_context_failed")
		logMock.Output.Reset()
	})
}

func getCreateExpenseRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"amount":893,"name":"Jordan shopping","period":"2023-5"}`,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}