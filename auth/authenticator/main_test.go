package main

import (
	"encoding/json"
	"github.com/JoelD7/money/api/storage"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func init() {
	storage.InitDynamoMock()
}

func TestLoginHandler(t *testing.T) {
	c := require.New(t)

	expectKey, result := expectedGetItem()
	storage.DynamoMock.ExpectGetItem().ToTable("person").WithKeys(expectKey).WillReturns(result)

	body := Credentials{
		Email:    "test@gmail.com",
		Password: "1234",
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	response, err := logInHandler(request)
	c.Equal(http.StatusOK, response.StatusCode)

}

func TestLoginHandlerFailed(t *testing.T) {
	c := require.New(t)

	expectKey, result := expectedGetItem()
	storage.DynamoMock.ExpectGetItem().ToTable("person").WithKeys(expectKey).WillReturns(result)

	body := Credentials{
		Email:    "test@gmail.com",
		Password: "random",
	}

	jsonBody, err := bodyToJSONString(body)
	c.Nil(err)

	request := &events.APIGatewayProxyRequest{Body: jsonBody}

	request.Body = jsonBody
	response, err := logInHandler(request)
	c.Equal(http.StatusBadRequest, response.StatusCode)
	c.Equal(errWrongCredentials.Error(), response.Body)

	storage.DynamoMock.ExpectGetItem().ToTable("person").WithKeys(expectKey).WillReturns(result)
}

func expectedGetItem() (map[string]*dynamodb.AttributeValue, dynamodb.GetItemOutput) {
	expectKey := map[string]*dynamodb.AttributeValue{
		"email": {
			S: aws.String("test@gmail.com"),
		},
	}

	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String("test@gmail.com"),
			},
			"password": {
				S: aws.String("$2a$10$.THF8QG33va8JTSIBz3lPuULaO6NiDb6yRmew63OtzujhVHbnZMFe"),
			},
			"full_name": {
				S: aws.String("Joel"),
			},
		},
	}

	return expectKey, result
}

func bodyToJSONString(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
