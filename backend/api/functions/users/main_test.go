package main

import (
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

func TestHandler(t *testing.T) {
	c := require.New(t)

	users.InitDynamoMock()
	income.InitDynamoMock()
	expenses.InitDynamoMock()

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{},
		PathParameters: map[string]string{
			"user-id": "test@gmail.com",
		},
	}

	response, err := handler(apigwRequest)
	c.Nil(err)
	c.NotNil(response)
	c.Equal(http.StatusOK, response.StatusCode)
	c.Contains(response.Body, `"Remainder":8670`)
}
