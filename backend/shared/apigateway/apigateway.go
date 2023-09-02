package apigateway

import (
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-lambda-go/events"
)

var (
	origin = env.GetString("CORS_ORIGIN", "*")
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// NewErrorResponse returns an error response
func NewErrorResponse(err error) *Response {
	var knownError *Error
	if errors.As(err, &knownError) {
		return NewJSONResponse(knownError.HTTPCode, knownError)
	}

	return NewJSONResponse(ErrInternalError.HTTPCode, ErrInternalError.Message)
}

// NewJSONResponse creates a new JSON response given a serializable `v`
func NewJSONResponse(statusCode int, v interface{}) *Response {
	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": origin,
		"Cache-Control":               "no-store",
		"Pragma":                      "no-cache",
		"Strict-Transport-Security":   "max-age=63072000; includeSubdomains; preload",
	}

	if origin != "*" {
		headers["Access-Control-Allow-Credentials"] = "true"
	}

	strData, ok := v.(string)
	if ok {
		return &Response{
			StatusCode: statusCode,
			Body:       strData,
			Headers:    headers,
		}
	}

	data, err := json.Marshal(v)
	if err != nil {
		return NewErrorResponse(errors.New("failed to marshal response"))
	}

	return &Response{
		StatusCode: statusCode,
		Body:       string(data),
		Headers:    headers,
	}
}
