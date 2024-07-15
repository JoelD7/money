package handlers

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreateSavingHandler(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	userMock := users.NewDynamoMock()
	savingsMock := savings.NewMock()
	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	dummyUser := users.GetDummyUser()
	dummyUser.Username = "username"

	err := userMock.CreateUser(ctx, dummyUser)
	c.NoError(err)

	req := &createSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
		userRepo:    userMock,
		periodRepo:  periodMock,
	}

	apigwRequest := getDummyRequest(dummyUser.Username)

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)

	userSavings, _, err := savingsMock.GetSavings(ctx, dummyUser.Username, "", 0)
	c.NoError(err)
	c.Len(userSavings, 1)
	c.Equal(dummyUser.Username, userSavings[0].Username)
}

func TestCreateSavingHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	userMock := users.NewDynamoMock()
	periodMock := period.NewDynamoMock()
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &createSavingRequest{
		log:         logMock,
		userRepo:    userMock,
		savingsRepo: savingsMock,
		periodRepo:  periodMock,
	}

	apigwRequest := getDummyRequest("")

	t.Run("Invalid request body - not JSON", func(t *testing.T) {
		apigwRequest = getDummyRequest("")
		apigwRequest.Body = "{"
		defer func() { apigwRequest = getDummyRequest("") }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Equal(http.StatusBadRequest, response.StatusCode)
		logMock.Output.Reset()
	})

	t.Run("Create saving failed", func(t *testing.T) {
		dummyError := errors.New("dummy error")

		savingsMock.ActivateForceFailure(dummyError)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Saving without amount", func(t *testing.T) {
		apigwRequest = getDummyRequest("")
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"test@gmail.com","period":"2020-01"}`
		defer func() { apigwRequest = getDummyRequest("") }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrMissingAmount.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Saving with invalid amount", func(t *testing.T) {
		apigwRequest = getDummyRequest("")
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"test@gmail.com","amount":-250,"period":"2020-01"}`
		defer func() { apigwRequest = getDummyRequest("") }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrInvalidSavingAmount.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing period", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"test@gmail.com","amount":250}`
		defer func() { apigwRequest = getDummyRequest("") }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrMissingPeriod.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})
}

func getDummyRequest(username string) *apigateway.Request {
	defaultUsername := "test@gmail.com"

	if username != "" {
		defaultUsername = username
	}

	return &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": defaultUsername,
			},
		},
		Body: `{"saving_goal_id":"SVG123","amount":250,"period":"2020-01"}`,
	}
}
