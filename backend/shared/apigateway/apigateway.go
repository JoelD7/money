package apigateway

import (
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"strings"
)

var (
	allowedOrigins    = env.GetString("CORS_ORIGIN", "")
	allowedOriginsMap = map[string]struct{}{}

	responseByErrors = map[error]Error{
		models.ErrUserNotFound:                   {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrIncomeNotFound:                 {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrExpenseNotFound:                {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrExpensesNotFound:               {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrCategoriesNotFound:             {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrSavingsNotFound:                {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrCategoryNotFound:               {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrInvalidAmount:                  {HTTPCode: http.StatusBadRequest, Message: "Invalid amount"},
		models.ErrMissingUsername:                {HTTPCode: http.StatusBadRequest, Message: "Missing username"},
		models.ErrInvalidEmail:                   {HTTPCode: http.StatusBadRequest, Message: "Invalid email"},
		models.ErrInvalidRequestBody:             {HTTPCode: http.StatusBadRequest, Message: "Invalid request body"},
		models.ErrMissingSavingID:                {HTTPCode: http.StatusBadRequest, Message: "Missing saving id"},
		models.ErrUpdateSavingNotFound:           {HTTPCode: http.StatusNotFound, Message: "The saving you are trying to update does not exist"},
		models.ErrDeleteSavingNotFound:           {HTTPCode: http.StatusNotFound, Message: "The saving you are trying to delete does not exist"},
		models.ErrInvalidPageSize:                {HTTPCode: http.StatusBadRequest, Message: "Invalid page size"},
		models.ErrInvalidStartKey:                {HTTPCode: http.StatusBadRequest, Message: "Invalid start key"},
		models.ErrMissingUsername:                {HTTPCode: http.StatusBadRequest, Message: "Missing username"},
		models.ErrMissingPassword:                {HTTPCode: http.StatusBadRequest, Message: "Missing password"},
		models.ErrInvalidToken:                   {HTTPCode: http.StatusUnauthorized, Message: "Invalid token"},
		models.ErrMalformedToken:                 {HTTPCode: http.StatusUnauthorized, Message: "Invalid token"},
		models.ErrExistingUser:                   {HTTPCode: http.StatusBadRequest, Message: "This account already exists"},
		models.ErrWrongCredentials:               {HTTPCode: http.StatusBadRequest, Message: "The email or password are incorrect"},
		models.ErrMissingCategoryName:            {HTTPCode: http.StatusBadRequest, Message: "Missing category name"},
		models.ErrInvalidHexColor:                {HTTPCode: http.StatusBadRequest, Message: "Invalid hex color"},
		models.ErrMissingCategoryColor:           {HTTPCode: http.StatusBadRequest, Message: "Missing category color"},
		models.ErrMissingCategoryBudget:          {HTTPCode: http.StatusBadRequest, Message: "Missing category budget"},
		models.ErrInvalidBudget:                  {HTTPCode: http.StatusBadRequest, Message: "Invalid budget"},
		models.ErrCategoryNameAlreadyExists:      {HTTPCode: http.StatusBadRequest, Message: "Category name already exists"},
		models.ErrMissingAmount:                  {HTTPCode: http.StatusBadRequest, Message: "Missing amount"},
		models.ErrInvalidSavingAmount:            {HTTPCode: http.StatusBadRequest, Message: "Invalid amount"},
		models.ErrSavingNotFound:                 {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrSavingGoalNotFound:             {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrNoUsernameInContext:            {HTTPCode: http.StatusBadRequest, Message: "Couldn't identify the user. Check if your Bearer token header is correct"},
		models.ErrMissingName:                    {HTTPCode: http.StatusBadRequest, Message: "Missing name"},
		models.ErrMissingExpenseID:               {HTTPCode: http.StatusBadRequest, Message: "Missing expense id"},
		models.ErrPeriodNotFound:                 {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrPeriodsNotFound:                {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrInvalidPeriod:                  {HTTPCode: http.StatusBadRequest, Message: "Invalid period"},
		models.ErrMissingPeriodDates:             {HTTPCode: http.StatusBadRequest, Message: "Missing period dates. A period should have a start_date and end_date"},
		models.ErrStartDateShouldBeBeforeEndDate: {HTTPCode: http.StatusBadRequest, Message: "start_date should be before end_date"},
		models.ErrPeriodNameIsTaken:              {HTTPCode: http.StatusBadRequest, Message: "Period name is taken"},
		models.ErrUpdatePeriodNotFound:           {HTTPCode: http.StatusNotFound, Message: "Not found"},
		models.ErrInvalidPeriodDate:              {HTTPCode: http.StatusBadRequest, Message: "Invalid period date"},
		models.ErrMissingPeriodID:                {HTTPCode: http.StatusBadRequest, Message: "Missing period id"},
		models.ErrMissingPeriod:                  {HTTPCode: http.StatusBadRequest, Message: "Missing period"},
		models.ErrMissingPeriodName:              {HTTPCode: http.StatusBadRequest, Message: "Missing period name"},
		models.ErrMissingPeriodStartDate:         {HTTPCode: http.StatusBadRequest, Message: "Missing period start date"},
		models.ErrMissingPeriodCreatedDate:       {HTTPCode: http.StatusBadRequest, Message: "Missing period created date"},
		models.ErrMissingPeriodUpdatedDate:       {HTTPCode: http.StatusBadRequest, Message: "Missing period updated date"},
		models.ErrExistingIncome:                 {HTTPCode: http.StatusBadRequest, Message: "This income already exists"},
		models.ErrMissingIncomeID:                {HTTPCode: http.StatusBadRequest, Message: "Missing income id"},
		models.ErrNoMoreItemsToBeRetrieved:       {HTTPCode: http.StatusNoContent, Message: "No more items to be retrieved"},
	}
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Header struct {
	Key   string
	Value string
}

func init() {
	origins := strings.Split(allowedOrigins, ";")
	for _, origin := range origins {
		allowedOriginsMap[origin] = struct{}{}
	}
}

// NewErrorResponse returns an error response
func (req *Request) NewErrorResponse(err error) *Response {
	var knownError *Error
	if errors.As(err, &knownError) {
		return req.NewJSONResponse(knownError.HTTPCode, knownError)
	}

	for mappedErr, responseErr := range responseByErrors {
		if errors.Is(err, mappedErr) {
			return req.NewJSONResponse(responseErr.HTTPCode, responseErr)
		}
	}

	return req.NewJSONResponse(ErrInternalError.HTTPCode, ErrInternalError)
}

// NewJSONResponse creates a new JSON response given a serializable `body`
func (req *Request) NewJSONResponse(statusCode int, body interface{}, headers ...Header) *Response {
	stdHeaders := map[string]string{
		"Content-Type":              "application/json",
		"Cache-Control":             "no-store",
		"Pragma":                    "no-cache",
		"Strict-Transport-Security": "max-age=63072000; includeSubdomains; preload",
	}

	origin := req.Headers["origin"]

	_, ok := allowedOriginsMap[origin]
	if ok {
		stdHeaders["Access-Control-Allow-Origin"] = origin
	}

	if allowedOrigins != "*" {
		stdHeaders["Access-Control-Allow-Credentials"] = "true"
	}

	for _, header := range headers {
		stdHeaders[header.Key] = header.Value
	}

	strData, ok := body.(string)
	if ok {
		return &Response{
			StatusCode: statusCode,
			Body:       strData,
			Headers:    stdHeaders,
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return req.NewErrorResponse(errors.New("failed to marshal response"))
	}

	return &Response{
		StatusCode: statusCode,
		Body:       string(data),
		Headers:    stdHeaders,
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
		"s_query_parameters":             paramsToString(req.QueryStringParameters),
		"s_headers":                      paramsToString(req.Headers),
		"s_multi_value_query_parameters": multiValueParamsToString(req.MultiValueQueryStringParameters),
		"s_path_parameters":              paramsToString(req.PathParameters),
		"o_authorizer":                   authorizer,
		"s_user_agent":                   req.Headers["User-Agent"],
		"s_content_type":                 req.Headers["Content-Type"],
		"s_method":                       req.HTTPMethod,
		"s_path":                         req.Path,
		"s_body":                         req.Body,
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

func multiValueParamsToString(params map[string][]string) string {
	var sb strings.Builder

	for param, values := range params {
		sb.WriteString(param)
		sb.WriteString("=")
		sb.WriteString(strings.Join(values, ","))
		sb.WriteString(" ")
	}

	return sb.String()
}

func GetUsernameFromContext(req *Request) (string, error) {
	username, ok := req.RequestContext.Authorizer["username"].(string)
	if !ok || username == "" {
		return "", models.ErrNoUsernameInContext
	}

	return username, nil
}
