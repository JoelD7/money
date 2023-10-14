package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetExpensesSuccess(t *testing.T) {
	c := require.New(t)

	dynamoClient := initDynamoClient()

	usersMock := users.NewDynamoRepository(dynamoClient)
	expensesMock := expenses.NewDynamoRepository(dynamoClient)
	//usersMock := users.NewDynamoMock()
	//expensesMock := expenses.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	request := &getExpensesRequest{
		log:          logMock,
		expensesRepo: expensesMock,
		userRepo:     usersMock,
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
}

func getGetExpensesRequest() *apigateway.Request {
	return &apigateway.Request{
		QueryStringParameters: map[string]string{
			"category": "test",
			"period":   "2023-7",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
