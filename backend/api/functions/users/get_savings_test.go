package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetSavingsHandler(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &getSavingsRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := getDummyAPIGatewayRequest()

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestGetSavingsHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &getSavingsRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := getDummyAPIGatewayRequest()

	t.Run("User savings not found", func(t *testing.T) {
		savingsMock.ActivateForceFailure(models.ErrSavingsNotFound)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(logMock.Output.String(), "savings_fetch_failed")
		c.Contains(logMock.Output.String(), models.ErrSavingsNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Invalid email", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{
			"email": "12",
		}
		defer func() { apigwRequest = getDummyAPIGatewayRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "invalid_email_detected")
		c.Contains(logMock.Output.String(), models.ErrInvalidEmail.Error())
		logMock.Output.Reset()
	})

	t.Run("Email not in authorizer context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{}
		defer func() { apigwRequest = getDummyAPIGatewayRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_user_email_from_context_failed")
		c.Contains(logMock.Output.String(), errNoUserEmailInContext.Error())
		logMock.Output.Reset()
	})
}

func getDummyAPIGatewayRequest() *apigateway.Request {
	return &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"email": "test@gmail.com",
			},
		},
	}
}