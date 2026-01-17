package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetPeriodsHandlerSuccess(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	periodMock := period.NewDynamoMock()

	request := &getPeriodsRequest{

		periodRepo: periodMock,
	}

	apigwRequest := getPeriodsAPIGatewayRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestGetPeriodsHandlerFailed(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	periodMock := period.NewDynamoMock()

	request := &getPeriodsRequest{

		periodRepo: periodMock,
	}

	apigwRequest := getPeriodsAPIGatewayRequest()

	t.Run("Username not in context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = nil
		defer func() { apigwRequest = getPeriodsAPIGatewayRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Invalid username", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer["username"] = "123"
		defer func() { apigwRequest = getPeriodsAPIGatewayRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Invalid page size parameter", func(t *testing.T) {
		apigwRequest.QueryStringParameters["page_size"] = "abc"
		defer func() { apigwRequest = getPeriodsAPIGatewayRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Periods not found", func(t *testing.T) {
		periodMock.ActivateForceFailure(models.ErrPeriodsNotFound)
		defer periodMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
	})
}

func getPeriodsAPIGatewayRequest() *apigateway.Request {
	return &apigateway.Request{
		QueryStringParameters: map[string]string{
			"page_size": "10",
			"start_key": "eyJlbWFpbCI6eyJWYWx1ZSI6InRlc3RAZ21haWwuY29tIn0sInNhdmluZ19pZCI6eyJWYWx1ZSI6IlNWMTU5In19",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
