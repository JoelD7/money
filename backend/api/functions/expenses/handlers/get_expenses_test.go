package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/expenses"
	periodRepo "github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetExpensesSuccess(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	expensesMock := expenses.NewDynamoMock()

	ctx := context.Background()

	request := &GetExpensesRequest{
		ExpensesRepo: expensesMock,
		UserRepo:     usersMock,
		PeriodRepo:   periodRepo.NewDynamoMock(),
	}

	apigwRequest := getGetExpensesRequest()

	t.Run("Query By period", func(t *testing.T) {
		apigwRequest.QueryStringParameters = map[string]string{
			"period": "2023-7",
		}
		defer func() { apigwRequest = getGetExpensesRequest() }()

		err := request.prepareRequest(apigwRequest)
		c.NoError(err)

		response, err := request.routeToHandlers(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusOK, response.StatusCode)
	})

	t.Run("Query by category", func(t *testing.T) {
		apigwRequest.MultiValueQueryStringParameters = map[string][]string{
			"category": {"CTGiBScOP3V16LYBjdIStP9"},
		}
		defer func() { apigwRequest = getGetExpensesRequest() }()

		err := request.prepareRequest(apigwRequest)
		c.NoError(err)

		response, err := request.routeToHandlers(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusOK, response.StatusCode)
	})
}

func getGetExpensesRequest() *apigateway.Request {
	return &apigateway.Request{
		QueryStringParameters: map[string]string{
			"category":  "test",
			"period":    "2023-7",
			"page_size": "5",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
