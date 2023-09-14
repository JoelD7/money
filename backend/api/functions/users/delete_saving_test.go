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

func TestDeleteHandlerSuccess(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &deleteSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyDeleteRequestBody(),
	}

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestDeleteHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &deleteSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyDeleteRequestBody(),
	}

	t.Run("Invalid request body - not JSON", func(t *testing.T) {
		apigwRequest.Body = "{"
		defer func() { apigwRequest.Body = getDummyDeleteRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(logMock.Output.String(), "request_body_unmarshal_failed")
		c.Equal(http.StatusBadRequest, response.StatusCode)
		logMock.Output.Reset()
	})

	t.Run("Saving with invalid email", func(t *testing.T) {
		apigwRequest.Body = `{"saving_id":"SVG123","email":"12","amount":250}`
		defer func() { apigwRequest.Body = getDummyDeleteRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrInvalidEmail.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Saving without email", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","amount":250}`
		defer func() { apigwRequest.Body = getDummyDeleteRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrMissingEmail.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("No saving ID", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","email":"test@gmail.com","amount":250}`
		defer func() { apigwRequest.Body = getDummyDeleteRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Equal(models.ErrMissingSavingID.Error(), response.Body)
		c.Contains(logMock.Output.String(), "delete_saving_failed")
	})

	t.Run("Saving doesn't exist", func(t *testing.T) {
		e := &mockRequestFailure{}

		savingsMock.ActivateForceFailure(e)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(response.Body, models.ErrDeleteSavingNotFound.Error())
	})
}

func getDummyDeleteRequestBody() string {
	return `{"saving_id":"SV123","email":"test@gmail.com"}`
}
