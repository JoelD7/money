package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateSaving(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &updateSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyUpdateRequestBody(),
	}

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestUpdateSavingHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &updateSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyUpdateRequestBody(),
	}

	t.Run("No email", func(t *testing.T) {
		apigwRequest.Body = `{"saving_id":"SV123","saving_goal_id":"SVG123","amount":250}`
		defer func() { apigwRequest.Body = getDummyUpdateRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Equal(models.ErrMissingEmail.Error(), response.Body)
		c.Contains(logMock.Output.String(), "update_saving_failed")
	})

	t.Run("Invalid amount", func(t *testing.T) {
		apigwRequest.Body = `{"saving_id":"SV123","saving_goal_id":"SVG123","email":"test@gmail.com","amount":0}`
		defer func() { apigwRequest.Body = getDummyUpdateRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Equal(models.ErrInvalidAmount.Error(), response.Body)
		c.Contains(logMock.Output.String(), "update_saving_failed")
	})

	t.Run("No saving ID", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","email":"test@gmail.com","amount":250}`
		defer func() { apigwRequest.Body = getDummyUpdateRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Equal(models.ErrMissingSavingID.Error(), response.Body)
		c.Contains(logMock.Output.String(), "update_saving_failed")
	})

	t.Run("Saving doesn't exist", func(t *testing.T) {
		e := &mockRequestFailure{}

		savingsMock.ActivateForceFailure(e)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(response.Body, models.ErrUpdateSavingNotFound)
	})
}

type mockRequestFailure struct{}

func (e *mockRequestFailure) StatusCode() int   { return 0 }
func (e *mockRequestFailure) RequestID() string { return "" }
func (e *mockRequestFailure) Code() string      { return "ConditionalCheckFailedException" }
func (e *mockRequestFailure) Message() string   { return "" }
func (e *mockRequestFailure) OrigErr() error    { return nil }
func (e *mockRequestFailure) Error() string     { return "" }

func getDummyUpdateRequestBody() string {
	return `{"saving_id":"SV123","saving_goal_id":"SVG123","email":"test@gmail.com","amount":250}`
}
