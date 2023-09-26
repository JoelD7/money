package main

import (
	"context"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateCategoryHandlerFailed(t *testing.T) {
	c := require.New(t)

	var apigwRequest *apigateway.Request

	usersMock := users.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	req := &updateCategoryRequest{
		userRepo: usersMock,
		log:      logMock,
	}

	t.Run("No category id in path", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.PathParameters = map[string]string{
			"categoryID": "",
		}

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "get_category_id_from_path_failed")
		c.Contains(logMock.Output.String(), errNoCategoryIDInPath.Error())
	})

	t.Run("Unmarshal request body failed", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = "invalid"

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "unmarshal_request_body_failed")
	})
}

func getUpdateCategoryRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"category_id":"CTGrR7fO4ndmI0IthJ7Wg8f","category_name":"Entertainment","color":"#ff8733","budget":1000}`,
		PathParameters: map[string]string{
			"categoryID": "CTGad",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
