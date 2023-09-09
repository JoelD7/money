package main

import (
	"context"
	"errors"
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

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest.Body = "{"
		defer func() { apigwRequest.Body = getDummyRequestBody() }()

		response, err := req.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Create saving failed", func(t *testing.T) {
		dummyError := errors.New("dummy error")

		savingsMock.ActivateForceFailure(dummyError)
		defer savingsMock.DeactivateForceFailure()

		response, err := req.process(ctx, apigwRequest)
		c.ErrorIs(err, dummyError)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
	})
}

func getDummyRequestBody() string {
	return `{"saving_id":"SV123","saving_goal_id":"SVG123","email":"test@gmail.com","creation_date":"2023-09-08T20:05:41.02376-04:00","amount":250}`
}
