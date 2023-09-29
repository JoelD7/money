package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUpdateCategoryHandler(t *testing.T) {
	c := require.New(t)

	usersMock := users.NewDynamoMock()
	logMock := logger.NewLoggerMock(nil)
	ctx := context.Background()

	req := &updateCategoryRequest{
		userRepo: usersMock,
		log:      logMock,
	}

	apigwRequest := getUpdateCategoryRequest()
	response, err := req.process(ctx, apigwRequest)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)

	t.Run("Update success without sending name", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","color":"#ff8733","budget":1000}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusOK, response.StatusCode)
	})

	t.Run("Update success without sending budget", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":"Entertainment","color":"#ff8733"}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusOK, response.StatusCode)
	})

	t.Run("Update success without sending color", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":"Entertainment","budget":1000}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusOK, response.StatusCode)
	})
}

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
		apigwRequest.PathParameters = map[string]string{}

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
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
	})

	t.Run("Update category failed", func(t *testing.T) {
		usersMock.ActivateForceFailure(errors.New("dummy"))
		defer usersMock.DeactivateForceFailure()

		apigwRequest = getUpdateCategoryRequest()

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.Contains(logMock.Output.String(), "update_category_failed")
	})

	t.Run("Invalid budget", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":789,"color":"#ff8733","budget":-89}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
	})

	t.Run("Name should not be empty", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":"","color":"#ff8733","budget":1000}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "request_body_validation_failed")
	})

	t.Run("Invalid color", func(t *testing.T) {
		apigwRequest = getUpdateCategoryRequest()
		apigwRequest.Body = `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":"Streaming","color":"asdf","budget":1000}`

		response, err := req.process(ctx, apigwRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "update_category_failed")
		c.Contains(response.Body, models.ErrInvalidHexColor.Error())
	})

	t.Run("Category name already exists", func(t *testing.T) {
		apiGatewayRequest := getUpdateCategoryRequest()

		response, err := req.process(ctx, apiGatewayRequest)
		c.Nil(err)
		c.Equal(http.StatusBadRequest, response.StatusCode)
		c.Contains(logMock.Output.String(), "update_category_failed")
		c.Contains(logMock.Output.String(), models.ErrCategoryNameAlreadyExists.Error())
	})
}

func getUpdateCategoryRequest() *apigateway.Request {
	return &apigateway.Request{
		Body: `{"id":"CTGzJeEzCNz6HMTiPKwgPmj","name":"Entertainment","color":"#ff8733","budget":1000}`,
		PathParameters: map[string]string{
			"categoryID": "CTGzJeEzCNz6HMTiPKwgPmj",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": "test@gmail.com",
			},
		},
	}
}
