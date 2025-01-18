package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	logger.InitLogger(logger.ConsoleImplementation)

	os.Exit(m.Run())
}
func TestCreateCategoryHandlerFailed(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()

	ctx := context.Background()

	req := &createCategoryRequest{
		userRepo: usersMock,
	}

	t.Run("Invalid request body", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing name", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"color":"#ff8733","budget":1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing budget", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","color":"#ff8733"}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Invalid budget", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","color":"#ff8733","budget":-1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Missing color", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()
		apiGatewayRequest.Body = `{"name":"Entertainment","budget":1000}`

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Category name already exists", func(t *testing.T) {
		apiGatewayRequest := getCreateCategoryRequest()

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
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
