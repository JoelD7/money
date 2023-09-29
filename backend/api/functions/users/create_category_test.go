package main

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCreateCategoryHandlerFailed(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	req := &createCategoryRequest{
		userRepo: usersMock,
		log:      logMock,
	}

	t.Run("Invalid request body", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
		c.Contains(logMock.Output.String(), models.ErrInvalidRequestBody.Error())
	})

	t.Run("Missing name", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"color":"#ff8733","budget":1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingCategoryName.Error())
	})

	t.Run("Missing budget", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","color":"#ff8733"}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingCategoryBudget.Error())
	})

	t.Run("Invalid budget", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","color":"#ff8733","budget":-1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
		c.Contains(logMock.Output.String(), models.ErrInvalidBudget.Error())
	})

	t.Run("Missing color", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","budget":1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
		c.Contains(logMock.Output.String(), models.ErrMissingCategoryColor.Error())
	})

	t.Run("Category name already exists", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "create_category_failed")
		c.Contains(logMock.Output.String(), models.ErrCategoryNameAlreadyExists.Error())
	})
}

func getCreateCategoryRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"name":"Entertainment","color":"#ff8733","budget":1000}`,
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
