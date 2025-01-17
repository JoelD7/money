package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetSavingsHandler(t *testing.T) {
	c := require.New(t)

	savingsMock := savings.NewMock()
	savingGoalMock := savingoal.NewMock()
	ctx := context.Background()

	req := &getSavingsRequest{
		username: "test@gmail.com",

		savingsRepo:    savingsMock,
		savingGoalRepo: savingGoalMock,
	}

	apigwRequest := getDummyAPIGatewayRequest()

	response, err := req.getUserSavings(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestGetSavingsHandlerFailed(t *testing.T) {
	c := require.New(t)

	savingsMock := savings.NewMock()
	savingGoalMock := savingoal.NewMock()
	ctx := context.Background()

	req := &getSavingsRequest{
		username: "test@gmail.com",

		savingGoalRepo: savingGoalMock,
		savingsRepo:    savingsMock,
	}

	apigwRequest := getDummyAPIGatewayRequest()

	t.Run("User savings not found", func(t *testing.T) {
		savingsMock.ActivateForceFailure(models.ErrSavingsNotFound)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.getUserSavings(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
	})

	t.Run("Invalid email", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{
			"username": "12",
		}
		defer func() { apigwRequest = getDummyAPIGatewayRequest() }()

		err := req.prepareRequest(apigwRequest)
		c.ErrorIs(err, models.ErrInvalidEmail)
	})

	t.Run("Username not in authorizer context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{}
		defer func() { apigwRequest = getDummyAPIGatewayRequest() }()

		err := req.prepareRequest(apigwRequest)
		c.ErrorIs(err, models.ErrNoUsernameInContext)
	})
}

func getDummyAPIGatewayRequest() *apigateway.Request {
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
