package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestJoel(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	periodMock, err := period.NewDynamoRepository(dynamoClient, "periods", "unique_periods")

	request := &updatePeriodRequest{
		periodRepo: periodMock,
	}

	apigwRequest := getUpdatePeriodRequest()
	apigwRequest.PathParameters = map[string]string{
		"periodID": "PRD2v6cW6J4tEQfm1Upneuq",
	}
	apigwRequest.Body = `{"created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29T00:00:00Z","name":"2023-6","start_date":"2023-01-01T00:00:00Z","updated_date":"2023-01-11T00:00:00Z"}`

	_, err = request.process(ctx, apigwRequest)
	c.NoError(err)
}

func TestUpdatePeriodHandlerSuccess(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	request := &updatePeriodRequest{
		periodRepo: periodMock,
	}

	apigwRequest := getUpdatePeriodRequest()

	response, err := request.process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestUpdatePeriodHandlerFailed_Database(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	request := &updatePeriodRequest{
		periodRepo: periodMock,
	}

	apigwRequest := getUpdatePeriodRequest()

	t.Run("Period not found", func(t *testing.T) {
		periodMock.ActivateForceFailure(models.ErrUpdatePeriodNotFound)
		defer periodMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
	})
}

func TestUpdatePeriodHandlerFailed_InputValidation(t *testing.T) {
	c := require.New(t)

	periodRepo := period.NewDynamoMock()
	ctx := context.Background()

	request := &updatePeriodRequest{

		periodRepo: periodRepo,
	}

	apigwRequest := getUpdatePeriodRequest()

	t.Run("Missing period ID", func(t *testing.T) {
		apigwRequest.PathParameters = map[string]string{}
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("No username in context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{}
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing name", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing period dates", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)

		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","start_date":"2023-01-29","name":"","updated_date":"2023-01-11T00:00:00Z"}`
		response, err = request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Start date after end date", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","start_date":"2023-01-30","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing created date", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","end_date":"2023-01-29","name":"","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})
}

func getUpdatePeriodRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29T00:00:00Z","name":"2023-01","start_date":"2023-01-01T00:00:00Z","updated_date":"2023-01-11T00:00:00Z"}`,
		PathParameters: map[string]string{
			"periodID": "2023-1",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
