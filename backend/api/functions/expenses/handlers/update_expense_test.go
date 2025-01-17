package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateHandlerSuccess(t *testing.T) {
	c := require.New(t)

	expensesMock := expenses.NewDynamoMock()
	periodMock := period.NewDynamoMock()
	userMock := users.NewDynamoMock()

	ctx := context.Background()

	request := &updateExpenseRequest{

		expensesRepo: expensesMock,
		periodRepo:   periodMock,
		userRepo:     userMock,
	}

	apigwRequest := getUpdateExpenseRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestUpdateHandlerFailed(t *testing.T) {
	c := require.New(t)

	expensesMock := expenses.NewDynamoMock()
	periodMock := period.NewDynamoMock()
	userMock := users.NewDynamoMock()

	ctx := context.Background()

	request := &updateExpenseRequest{

		expensesRepo: expensesMock,
		periodRepo:   periodMock,
		userRepo:     userMock,
	}

	apigwRequest := getUpdateExpenseRequest()

	t.Run("Invalid period", func(t *testing.T) {
		apigwRequest.Body = `{"amount":892,"period":"2020-13"}`
		defer func() { apigwRequest = getUpdateExpenseRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})
}

func getUpdateExpenseRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"amount":892,"period":"2020-01"}`,
		PathParameters: map[string]string{
			"expenseID": "EX0H4ddQBWAkNFEUMdzLYY",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
