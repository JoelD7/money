package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreateSavingHandler(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &createSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyRequestBody(),
	}

	response, err := req.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func TestCreateSavingHandlerFailed(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	savingsMock := savings.NewMock()
	ctx := context.Background()

	req := &createSavingRequest{
		log:         logMock,
		savingsRepo: savingsMock,
	}

	apigwRequest := &apigateway.Request{
		Body: getDummyRequestBody(),
	}

	t.Run("Invalid request body - not JSON", func(t *testing.T) {
		apigwRequest.Body = "{"
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Equal(http.StatusBadRequest, response.StatusCode)
		logMock.Output.Reset()
	})

	t.Run("Missing username", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","amount":250}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

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

	t.Run("Saving with invalid email", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"12","amount":250}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrInvalidEmail.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Saving without username", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","amount":250}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrMissingUsername.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Saving without amount", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"test@gmail.com"}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrMissingSavingAmount.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Saving with invalid amount", func(t *testing.T) {
		apigwRequest.Body = `{"saving_goal_id":"SVG123","username":"test@gmail.com","amount":-250}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrInvalidSavingAmount.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Empty request body", func(t *testing.T) {
		apigwRequest.Body = `{}`
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(models.ErrInvalidRequestBody.Error(), response.Body)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})
}

func getDummyRequestBody() string {
	return `{"saving_goal_id":"SVG123","username":"test@gmail.com","amount":250}`
}
