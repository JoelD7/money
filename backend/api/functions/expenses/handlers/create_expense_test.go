package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	logger.InitLogger(logger.ConsoleImplementation)

	os.Exit(m.Run())
}

func TestHandlerSuccess(t *testing.T) {
	c := require.New(t)

	userMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	periodMock := period.NewDynamoMock()
	cacheMock := cache.NewRedisCacheMock()

	ctx := context.Background()

	request := &createExpenseRequest{
		userRepo:         userMock,
		expensesRepo:     expensesMock,
		periodRepo:       periodMock,
		idempotenceCache: cacheMock,
	}

	apigwRequest := getCreateExpenseRequest(periodMock)

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func TestHandlerFailure(t *testing.T) {
	c := require.New(t)

	userMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()
	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	request := &createExpenseRequest{
		idempotenceCache: cache.NewRedisCacheMock(),
		userRepo:         userMock,
		expensesRepo:     expensesMock,
		periodRepo:       periodMock,
	}

	apigwRequest := getCreateExpenseRequest(periodMock)

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest.Body = "invalid"
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Invalid email", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer["username"] = "invalid"
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing name", func(t *testing.T) {
		apigwRequest.Body = `{"amount":893,"period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, "Missing name")
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing amount", func(t *testing.T) {
		apigwRequest.Body = `{"name":"Jordan shopping","period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, "Missing amount")
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		apigwRequest.Body = `{"amount":0,"name":"Jordan shopping","period":"2023-5"}`
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, "Invalid amount")
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Create expense failed", func(t *testing.T) {
		expensesMock.ActivateForceFailure(errors.New("dummy"))
		defer expensesMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Get Username from context failed", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = nil
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Missing period", func(t *testing.T) {
		apigwRequest.Body = `{"amount":893,"name":"Jordan shopping"}`
		defer func() { apigwRequest = getCreateExpenseRequest(periodMock) }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(response.Body, "Missing period")
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})
}

func getCreateExpenseRequest(periodMock *period.DynamoMock) *apigateway.Request {
	body := fmt.Sprintf(`{"amount":893,"name":"Jordan shopping","period_id":"%s"}`, periodMock.GetDefaultPeriod().ID)

	return &apigateway.Request{
		Body: body,
		Headers: map[string]string{
			"Idempotency-Key": "123",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
