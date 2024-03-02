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
	allowedOrigins    = env.GetString("CORS_ORIGIN", "https://localhost")
	allowedOriginsMap = map[string]struct{}{}

	responseByErrors = map[error]Error{
		models.ErrUserNotFound:                   {HTTPCode: http.StatusNotFound, Message: models.ErrUserNotFound.Error()},
		models.ErrIncomeNotFound:                 {HTTPCode: http.StatusNotFound, Message: models.ErrIncomeNotFound.Error()},
		models.ErrExpenseNotFound:                {HTTPCode: http.StatusNotFound, Message: models.ErrExpenseNotFound.Error()},
		models.ErrExpensesNotFound:               {HTTPCode: http.StatusNotFound, Message: models.ErrExpensesNotFound.Error()},
		models.ErrCategoriesNotFound:             {HTTPCode: http.StatusNotFound, Message: models.ErrCategoriesNotFound.Error()},
		models.ErrSavingsNotFound:                {HTTPCode: http.StatusNotFound, Message: models.ErrSavingsNotFound.Error()},
		models.ErrCategoryNotFound:               {HTTPCode: http.StatusNotFound, Message: models.ErrCategoryNotFound.Error()},
		models.ErrInvalidAmount:                  {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidAmount.Error()},
		models.ErrMissingUsername:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingUsername.Error()},
		models.ErrInvalidEmail:                   {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidEmail.Error()},
		models.ErrInvalidRequestBody:             {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidRequestBody.Error()},
		models.ErrMissingSavingID:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingSavingID.Error()},
		models.ErrUpdateSavingNotFound:           {HTTPCode: http.StatusNotFound, Message: models.ErrUpdateSavingNotFound.Error()},
		models.ErrDeleteSavingNotFound:           {HTTPCode: http.StatusNotFound, Message: models.ErrDeleteSavingNotFound.Error()},
		models.ErrInvalidPageSize:                {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidPageSize.Error()},
		models.ErrInvalidStartKey:                {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidStartKey.Error()},
		models.ErrMissingUsername:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingUsername.Error()},
		models.ErrMissingPassword:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPassword.Error()},
		models.ErrInvalidToken:                   {HTTPCode: http.StatusUnauthorized, Message: models.ErrInvalidToken.Error()},
		models.ErrMalformedToken:                 {HTTPCode: http.StatusUnauthorized, Message: models.ErrMalformedToken.Error()},
		models.ErrExistingUser:                   {HTTPCode: http.StatusBadRequest, Message: models.ErrExistingUser.Error()},
		models.ErrWrongCredentials:               {HTTPCode: http.StatusBadRequest, Message: models.ErrWrongCredentials.Error()},
		models.ErrMissingCategoryName:            {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingCategoryName.Error()},
		models.ErrInvalidHexColor:                {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidHexColor.Error()},
		models.ErrMissingCategoryColor:           {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingCategoryColor.Error()},
		models.ErrMissingCategoryBudget:          {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingCategoryBudget.Error()},
		models.ErrInvalidBudget:                  {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidBudget.Error()},
		models.ErrCategoryNameAlreadyExists:      {HTTPCode: http.StatusBadRequest, Message: models.ErrCategoryNameAlreadyExists.Error()},
		models.ErrMissingAmount:                  {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingAmount.Error()},
		models.ErrInvalidSavingAmount:            {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidSavingAmount.Error()},
		models.ErrSavingNotFound:                 {HTTPCode: http.StatusNotFound, Message: models.ErrSavingNotFound.Error()},
		models.ErrSavingGoalNotFound:             {HTTPCode: http.StatusNotFound, Message: models.ErrSavingGoalNotFound.Error()},
		models.ErrNoUsernameInContext:            {HTTPCode: http.StatusBadRequest, Message: models.ErrNoUsernameInContext.Error()},
		models.ErrMissingName:                    {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingName.Error()},
		models.ErrMissingExpenseID:               {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingExpenseID.Error()},
		models.ErrPeriodNotFound:                 {HTTPCode: http.StatusNotFound, Message: models.ErrPeriodNotFound.Error()},
		models.ErrPeriodsNotFound:                {HTTPCode: http.StatusNotFound, Message: models.ErrPeriodsNotFound.Error()},
		models.ErrInvalidPeriod:                  {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidPeriod.Error()},
		models.ErrMissingPeriodDates:             {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodDates.Error()},
		models.ErrStartDateShouldBeBeforeEndDate: {HTTPCode: http.StatusBadRequest, Message: models.ErrStartDateShouldBeBeforeEndDate.Error()},
		models.ErrPeriodNameIsTaken:              {HTTPCode: http.StatusBadRequest, Message: models.ErrPeriodNameIsTaken.Error()},
		models.ErrUpdatePeriodNotFound:           {HTTPCode: http.StatusNotFound, Message: models.ErrUpdatePeriodNotFound.Error()},
		models.ErrInvalidPeriodDate:              {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidPeriodDate.Error()},
		models.ErrMissingPeriodID:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodID.Error()},
		models.ErrMissingPeriod:                  {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriod.Error()},
		models.ErrMissingPeriodName:              {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodName.Error()},
		models.ErrMissingPeriodStartDate:         {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodStartDate.Error()},
		models.ErrMissingPeriodCreatedDate:       {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodCreatedDate.Error()},
		models.ErrMissingPeriodUpdatedDate:       {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPeriodUpdatedDate.Error()},
		models.ErrExistingIncome:                 {HTTPCode: http.StatusBadRequest, Message: models.ErrExistingIncome.Error()},
		models.ErrMissingIncomeID:                {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingIncomeID.Error()},
		models.ErrNoMoreItemsToBeRetrieved:       {HTTPCode: http.StatusNoContent, Message: models.ErrNoMoreItemsToBeRetrieved.Error()},
	}
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

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
		return req.NewJSONResponse(knownError.HTTPCode, knownError.Message)
	}

	for mappedErr, responseErr := range responseByErrors {
		if errors.Is(err, mappedErr) {
			return req.NewJSONResponse(responseErr.HTTPCode, responseErr.Message)
		}
	}

	return req.NewJSONResponse(ErrInternalError.HTTPCode, ErrInternalError.Message)
}

// NewJSONResponse creates a new JSON response given a serializable `body`
func (req *Request) NewJSONResponse(statusCode int, body interface{}) *Response {
	origin := req.Headers["Origin"]

	allowOriginHeader := ""
	if _, ok := allowedOriginsMap[origin]; ok {
		allowOriginHeader = origin
	}

	headers := map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": allowOriginHeader,
		"Cache-Control":               "no-store",
		"Pragma":                      "no-cache",
		"Strict-Transport-Security":   "max-age=63072000; includeSubdomains; preload",
	}

	if allowedOrigins != "*" {
		headers["Access-Control-Allow-Credentials"] = "true"
		headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS, PUT, DELETE"
		headers["Access-Control-Allow-Headers"] = "Content-Type"
	}

	strData, ok := body.(string)
	if ok {
		return &Response{
			StatusCode: statusCode,
			Body:       strData,
			Headers:    headers,
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return req.NewErrorResponse(errors.New("failed to marshal response"))
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
		"s_query_parameters":             paramsToString(req.QueryStringParameters),
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
