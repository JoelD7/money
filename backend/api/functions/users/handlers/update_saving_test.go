package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateSaving(t *testing.T) {
	c := require.New(t)

	savingsMock := savings.NewMock()
	savingGoalMock := savingoal.NewMock()
	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	req := &updateSavingRequest{
		savingsRepo:    savingsMock,
		periodRepo:     periodMock,
		savingGoalRepo: savingGoalMock,
	}

	apigwRequest := getDummyUpdateRequest()

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestUpdateSavingHandlerFailed(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &updateSavingRequest{
		savingsRepo:    savingsMock,
		savingGoalRepo: savingoal.NewMock(),
		periodRepo:     periodMock,
	}

	apigwRequest := getDummyUpdateRequest()

	t.Run("Invalid email", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer["username"] = "test"
		defer func() { apigwRequest = getDummyUpdateRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, "Invalid email")
	})

	t.Run("Invalid amount", func(t *testing.T) {
		apigwRequest.Body = `{"saving_id":"SV123","saving_goal_id":"SVG123","amount":0}`
		defer func() { apigwRequest = getDummyUpdateRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, "Invalid amount")
	})

	t.Run("No saving ID", func(t *testing.T) {
		apigwRequest.PathParameters["savingID"] = ""
		defer func() { apigwRequest = getDummyUpdateRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, "Missing saving id")
	})

	t.Run("Saving doesn't exist", func(t *testing.T) {
		e := &mockRequestFailure{}

		savingsMock.ActivateForceFailure(e)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(response.Body, "The saving you are trying to update does not exist")
	})

}

type mockRequestFailure struct{}

func (e *mockRequestFailure) StatusCode() int   { return 0 }
func (e *mockRequestFailure) RequestID() string { return "" }
func (e *mockRequestFailure) Code() string      { return "ConditionalCheckFailedException" }
func (e *mockRequestFailure) Message() string   { return "" }
func (e *mockRequestFailure) OrigErr() error    { return nil }
func (e *mockRequestFailure) Error() string     { return "ConditionalCheckFailedException" }

func getDummyUpdateRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"saving_goal_id":"SVG123","username":"test@gmail.com","amount":250,"period":"2020-01"}`,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
		PathParameters: map[string]string{
			"savingID": "SV123",
		},
	}
}
