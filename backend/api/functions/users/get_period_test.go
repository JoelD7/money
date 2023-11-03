package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetPeriodHandlerSuccess(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()
	periodMock := period.NewDynamoMock()

	request := &getPeriodRequest{
		log:        logMock,
		periodRepo: periodMock,
	}

	apigwRequest := getPeriodAPIGatewayRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestGetPeriodHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()
	periodMock := period.NewDynamoMock()

	request := &getPeriodRequest{
		log:        logMock,
		periodRepo: periodMock,
	}

	apigwRequest := getPeriodAPIGatewayRequest()

	t.Run("Missing period ID", func(t *testing.T) {
		apigwRequest.PathParameters = nil
		defer func() { apigwRequest = getPeriodAPIGatewayRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "missing_period_id")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodID.Error())
	})

	t.Run("Username not in context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = nil
		defer func() { apigwRequest = getPeriodAPIGatewayRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_username_from_context_failed")
		c.Contains(logMock.Output.String(), models.ErrNoUsernameInContext.Error())
	})
}

func getPeriodAPIGatewayRequest() *apigateway.Request {
	return &apigateway.Request{
		PathParameters: map[string]string{
			"periodID": "123",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
