package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreatePeriodSuccess(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()
	ctx := context.Background()

	request := &CreatePeriodRequest{

		PeriodRepo: periodMock,
	}

	apigwRequest := getCreatePeriodRequest()

	response, err := request.Process(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusCreated, response.StatusCode)
}

func TestCreatePeriodSuccessFailed(t *testing.T) {
	c := require.New(t)

	periodMock := period.NewDynamoMock()

	ctx := context.Background()

	request := &CreatePeriodRequest{
		PeriodRepo: periodMock,
	}

	apigwRequest := getCreatePeriodRequest()

	t.Run("Missing period name", func(t *testing.T) {
		apigwRequest.Body = `{"start_date":"2023-12-01","end_date":"2023-12-05"}`
		defer func() { apigwRequest = getCreatePeriodRequest() }()

		response, err := request.Process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, models.ErrMissingPeriodName.Error())
	})
}

func getCreatePeriodRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"start_date":"2023-12-01","end_date":"2023-12-05","name":"2023-2"}`,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
