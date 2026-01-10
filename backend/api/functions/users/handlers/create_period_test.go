package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
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
		IdempotenceCache:         cache.NewRedisCacheMock(),
		IncomePeriodCacheManager: cache.NewRedisCacheMock(),
		PeriodRepo:               periodMock,
		SavingGoalRepo:           savingoal.NewMock(),
		SavingsRepo:              savings.NewMock(),
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
		PeriodRepo:               periodMock,
		IncomePeriodCacheManager: cache.NewRedisCacheMock(),
		IdempotenceCache:         cache.NewRedisCacheMock(),
		SavingGoalRepo:           savingoal.NewMock(),
		SavingsRepo:              savings.NewMock(),
	}

	apigwRequest := getCreatePeriodRequest()

	t.Run("Missing period name", func(t *testing.T) {
		apigwRequest.Body = `{"start_date":"2023-12-01T15:04:05Z","end_date":"2023-12-05T15:04:05Z"}`
		defer func() { apigwRequest = getCreatePeriodRequest() }()

		response, err := request.Process(ctx, apigwRequest)
		c.NoError(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(response.Body, "Missing period name")
	})
}

func getCreatePeriodRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"start_date":"2023-12-01T15:04:05Z","end_date":"2023-12-05T15:04:05Z","name":"2023-2"}`,
		Headers: map[string]string{
			"Idempotency-Key": "123",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
