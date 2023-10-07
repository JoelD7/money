package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetSavingHandler(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	savingoalMock := savingoal.NewMock()
	ctx := context.Background()

	req := &getSavingRequest{
		log:            logMock,
		savingsRepo:    savingsMock,
		savingGoalRepo: savingoalMock,
	}

	apigwRequest := getSavingAPIGatewayRequest()

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestGetSavingHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	savingoalMock := savingoal.NewMock()
	ctx := context.Background()

	req := &getSavingRequest{
		log:            logMock,
		savingsRepo:    savingsMock,
		savingGoalRepo: savingoalMock,
	}

	apigwRequest := getSavingAPIGatewayRequest()

	t.Run("Saving not found", func(t *testing.T) {
		savingsMock.ActivateForceFailure(models.ErrSavingNotFound)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_saving_failed")
		c.Contains(logMock.Output.String(), models.ErrSavingNotFound.Error())
		logMock.Output.Reset()
	})

	t.Run("Missing savingID", func(t *testing.T) {
		apigwRequest.PathParameters = map[string]string{}
		defer func() { apigwRequest = getSavingAPIGatewayRequest() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "missing_saving_id")
		logMock.Output.Reset()
	})

	t.Run("Get saving goal name failed", func(t *testing.T) {
		savingoalMock.ActivateForceFailure(models.ErrSavingGoalNameSettingFailed)
		defer savingoalMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusOK, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_saving_goal_name_failed")
		logMock.Output.Reset()
	})
}

func getSavingAPIGatewayRequest() *apigateway.Request {
	savingID := "dummy"
	if len(savings.GetDummySavings()) > 0 {
		savingID = savings.GetDummySavings()[0].SavingID
	}

	return &apigateway.Request{
		PathParameters: map[string]string{
			"savingID": savingID,
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
