package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func Test(t *testing.T) {
	c := require.New(t)

	//dynamoClient := initDynamoClient()
	//periodRepo := period.NewDynamoRepository(dynamoClient)
	periodRepo := period.NewDynamoMock()

	err := periodRepo.UpdatePeriod(context.Background(), &models.Period{
		ID:        "2023-1",
		Username:  "test@gmail.com",
		Name:      aws.String(""),
		StartDate: models.ToPeriodTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   models.ToPeriodTime(time.Date(2023, 1, 29, 0, 0, 0, 0, time.UTC)),
	})
	c.Nil(err)
}

func TestUpdatePeriodHandlerFailed_Database(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	request := &updatePeriodRequest{
		log:        logMock,
		periodRepo: periodMock,
	}

	apigwRequest := getUpdatePeriodRequest()

	t.Run("Period not found", func(t *testing.T) {
		periodMock.ActivateForceFailure(models.ErrUpdatePeriodNotFound)
		defer periodMock.DeactivateForceFailure()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusNotFound, response.StatusCode)
		c.Contains(logMock.Output.String(), "update_period_failed")
		c.Contains(logMock.Output.String(), models.ErrUpdatePeriodNotFound.Error())
	})
}

func TestUpdatePeriodHandlerFailed_InputValidation(t *testing.T) {
	c := require.New(t)

	logMock := logger.NewLoggerMock(nil)
	periodRepo := period.NewDynamoMock()
	ctx := context.Background()

	request := &updatePeriodRequest{
		log:        logMock,
		periodRepo: periodRepo,
	}

	apigwRequest := getUpdatePeriodRequest()

	t.Run("Missing period ID", func(t *testing.T) {
		apigwRequest.PathParameters = map[string]string{}
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodID.Error())
		logMock.Output.Reset()
	})

	t.Run("No username in context", func(t *testing.T) {
		apigwRequest.RequestContext.Authorizer = map[string]interface{}{}
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrNoUsernameInContext.Error())
		logMock.Output.Reset()
	})

	t.Run("Invalid request body", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrInvalidRequestBody.Error())
		logMock.Output.Reset()
	})

	t.Run("Missing name", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodName.Error())
		logMock.Output.Reset()
	})

	t.Run("Missing period dates", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodDates.Error())
		logMock.Output.Reset()

		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","start_date":"2023-01-29","name":"","updated_date":"2023-01-11T00:00:00Z"}`
		response, err = request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodDates.Error())
		logMock.Output.Reset()
	})

	t.Run("Start date after end date", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","start_date":"2023-01-30","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrStartDateShouldBeBeforeEndDate.Error())
		logMock.Output.Reset()
	})

	t.Run("Missing created date", func(t *testing.T) {
		apigwRequest.Body = `{"period":"2023-1","end_date":"2023-01-29","name":"","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"}`
		defer func() { apigwRequest = getUpdatePeriodRequest() }()

		response, err := request.process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "validate_request_body_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingPeriodCreatedDate.Error())
		logMock.Output.Reset()
	})
}

func getUpdatePeriodRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"period":"2023-1","created_date":"2023-10-21T17:53:21.908187368Z","end_date":"2023-01-29","name":"","start_date":"2023-01-01","updated_date":"2023-01-11T00:00:00Z"}`,
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
