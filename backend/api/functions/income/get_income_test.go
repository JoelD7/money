package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetIncome(t *testing.T) {
	c := require.New(t)

	dynamoClient := initDynamoClient()

	request := &incomeGetRequest{
		incomeRepo: income.NewDynamoRepository(dynamoClient),
		log:        logger.NewLogger(),
	}

	apigwRequest := &apigateway.Request{
		PathParameters: map[string]string{
			"incomeID": "dummy",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}

	defer func() {
		err := request.log.Close()
		c.Nil(err)
	}()

	response, err := request.process(context.Background(), apigwRequest)
	c.Nil(err)
	c.NotNil(response)
}
