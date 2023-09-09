package apigateway

import (
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-lambda-go/events"
	"strings"
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

func (req *Request) LogName() string {
	return "http_request"
}

func (req *Request) LogProperties() map[string]interface{} {
	authorizer := map[string]interface{}{
		"s_event_id":        req.RequestContext.Authorizer["event_id"],
		"s_username":        req.RequestContext.Authorizer["username"],
		"s_client_id":       req.RequestContext.Authorizer["client_id"],
		"s_scope":           req.RequestContext.Authorizer["scope"],
		"s_api_key_version": req.RequestContext.Authorizer["version"],
		"b_is_internal":     req.RequestContext.Authorizer["is_internal"],
	}

	return map[string]interface{}{
		"s_query_parameters": paramsToString(req.QueryStringParameters),
		"s_path_parameters":  paramsToString(req.PathParameters),
		"o_authorizer":       authorizer,
		"s_user_agent":       req.Headers["User-Agent"],
		"s_content_type":     req.Headers["Content-Type"],
		"s_method":           req.HTTPMethod,
		"s_path":             req.Path,
		"s_body":             req.Body,
	}
}

func paramsToString(params map[string]string) string {
	var sb strings.Builder

	for param, value := range params {
		sb.WriteString(param)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteString(" ")
	}

	return sb.String()
}
