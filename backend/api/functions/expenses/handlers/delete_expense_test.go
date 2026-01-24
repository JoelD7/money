package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestDeleteHandlerSuccess(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()

	expensesMock := expenses.NewDynamoMock()

	request := &deleteExpenseRequest{
		expensesRepo: expensesMock,
	}

	apiRequest := getDeleteExpenseRequest()

	response, err := request.process(ctx, apiRequest)
	c.NoError(err)
	c.Equal(http.StatusNoContent, response.StatusCode)
}

func getDeleteExpenseRequest() *apigateway.Request {
	return &apigateway.Request{
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
