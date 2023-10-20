package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateHandlerSuccess(t *testing.T) {
	c := require.New(t)

	dynamoClient := initDynamoClient()

	//expensesMock := expenses.NewDynamoMock()
	expensesMock := expenses.NewDynamoRepository(dynamoClient)
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	request := &updateExpenseRequest{
		log:          logMock,
		expensesRepo: expensesMock,
	}

	apigwRequest := getUpdateExpenseRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func getUpdateExpenseRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"amount":892}`,
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
