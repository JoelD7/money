package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreatePeriodSuccess(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	request := &createPeriodRequest{
		log:        logMock,
		periodRepo: periodMock,
	}

	apigwRequest := getCreatePeriodRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func getCreatePeriodRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"start_date":"2023-10-28","end_date":"2023-10-29"}`,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
