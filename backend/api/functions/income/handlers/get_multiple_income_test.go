package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetIncomeByPeriod(t *testing.T) {
	c := require.New(t)

	apigwRequest := getDummyAPIGatewayRequest()
	apigwRequest.QueryStringParameters["period"] = "2023-7"

	ctx := context.Background()
	incomeMock := income.NewDynamoMock()

	request := &GetMultipleIncomeRequest{

		IncomeRepo: incomeMock,
	}

	err := request.prepareRequest(apigwRequest)
	c.NoError(err)

	response, err := request.RouteToHandlers(ctx, apigwRequest)
	c.NoError(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func getDummyAPIGatewayRequest() *apigateway.Request {
	return &apigateway.Request{
		QueryStringParameters: map[string]string{
			"page_size": "10",
			//"start_key": "eyJpbmNvbWVfaWQiOiJJTnlTQVU3bnN1TFJBbkpyVHZybGcwIiwicGVyaW9kX3VzZXIiOiIyMDIzLTc6dGVzdEBnbWFpbC5jb20ifQ==",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
