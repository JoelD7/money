package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/usecases"
	"strings"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Effect int

type HttpVerb int

type AuthorizerResponse struct {
	events.APIGatewayCustomAuthorizerResponse

	// The region where the API is deployed. By default this is set to '*'
	Region string

	// The AWS account id the policy will be generated for. This is used to create the method ARNs.
	AccountID string

	// The API Gateway API id. By default this is set to '*'
	APIID string

	// The name of the stage used in the policy. By default this is set to '*'
	Stage    string
	Resource string
}

const (
	Get HttpVerb = iota
	Post
	Put
	Delete
	Patch
	Head
	Options
	All
)

const (
	Allow Effect = iota
	Deny
)

var (
	kidSecretName = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	awsRegion     = env.GetString("REGION", "us-east-1")

	errUserNotAuthorized = errors.New("access to this user's data is forbidden")
)

type request struct {
	log            logger.LogAPI
	secretsManager secrets.SecretManager
	cacheRepo      cache.InvalidTokenManager
	startingTime   time.Time
	client         restclient.HttpClient
	err            error
}

func (req *request) init() {
	req.startingTime = time.Now()
	req.cacheRepo = cache.NewRedisCache()
	req.secretsManager = secrets.NewAWSSecretManager()
	req.client = restclient.New()
}

func (req *request) finish() {
	defer func() {
		err := req.log.Close()
		if err != nil {
			panic(err)
		}
	}()

	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	req := &request{
		log: logger.NewLogger(),
	}

	req.init()
	defer req.finish()

	return req.process(ctx, event)
}

func (req *request) process(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.ReplaceAll(event.AuthorizationToken, "Bearer ", "")

	verifyToken := usecases.NewTokenVerifier(req.client, req.log, req.secretsManager, req.cacheRepo)

	subject, err := verifyToken(ctx, token)
	if errors.Is(err, models.ErrUnauthorized) {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	if err != nil {
		return defaultDenyAllPolicy(event.MethodArn, err), nil
	}

	principalID := subject

	resp := NewAuthorizerResponse(event.MethodArn, principalID)

	err = req.verifyUserRequest(resp, event.MethodArn)
	if err != nil {
		req.log.Error("request_denied", err, []models.LoggerObject{req.getEventAsLoggerObject(event)})

		return defaultDenyAllPolicy(event.MethodArn, err), nil
	}

	resp.AllowAllMethods()

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func (req *request) getEventAsLoggerObject(event events.APIGatewayCustomAuthorizerRequest) models.LoggerObject {
	return req.log.MapToLoggerObject("authorizer_request", map[string]interface{}{
		"s_type":       event.Type,
		"s_method_arn": event.MethodArn,
	})
}

func (req *request) verifyUserRequest(resp *AuthorizerResponse, methodArn string) error {
	arnParts := strings.Split(methodArn, ":")
	apiGatewayArnParts := strings.Split(arnParts[5], "/")

	if len(apiGatewayArnParts) < 5 {
		return nil
	}

	userID := apiGatewayArnParts[4]

	if resp.PrincipalID != userID {
		return errUserNotAuthorized
	}

	return nil
}

func defaultDenyAllPolicy(methodArn string, err error) events.APIGatewayCustomAuthorizerResponse {
	resp := NewAuthorizerResponse(methodArn, "user")
	resp.DenyAllMethods()

	resp.APIGatewayCustomAuthorizerResponse.Context = map[string]interface{}{
		"stringKey": err.Error(),
	}

	return resp.APIGatewayCustomAuthorizerResponse
}

func (hv HttpVerb) String() string {
	switch hv {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case Patch:
		return "PATCH"
	case Head:
		return "HEAD"
	case Options:
		return "OPTIONS"
	case All:
		return "*"
	}
	return ""
}

func (e Effect) String() string {
	switch e {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	}
	return ""
}

func NewAuthorizerResponse(methodArn, principalID string) *AuthorizerResponse {
	tmp := strings.Split(methodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")

	return &AuthorizerResponse{
		APIGatewayCustomAuthorizerResponse: events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: principalID,
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
			},
			Context: map[string]interface{}{
				"username": principalID,
			},
		},
		Region:    awsRegion,
		AccountID: tmp[4],
		APIID:     apiGatewayArnTmp[0],
		Stage:     apiGatewayArnTmp[1],
	}
}

func (r *AuthorizerResponse) addMethod(effect Effect, verb HttpVerb, resource string) {
	resourceArn := "arn:aws:execute-api:" +
		r.Region + ":" +
		r.AccountID + ":" +
		r.APIID + "/" +
		r.Stage + "/" +
		verb.String() + "/" +
		strings.TrimLeft(resource, "/")

	s := events.IAMPolicyStatement{
		Effect:   effect.String(),
		Action:   []string{"execute-api:Invoke"},
		Resource: []string{resourceArn},
	}

	r.PolicyDocument.Statement = append(r.PolicyDocument.Statement, s)
}

func (r *AuthorizerResponse) AllowAllMethods() {
	r.addMethod(Allow, All, "*")
}

func (r *AuthorizerResponse) DenyAllMethods() {
	r.addMethod(Deny, All, "*")
}

func (r *AuthorizerResponse) AllowMethod(verb HttpVerb, resource string) {
	r.addMethod(Allow, verb, resource)
}

func (r *AuthorizerResponse) DenyMethod(verb HttpVerb, resource string) {
	r.addMethod(Deny, verb, resource)
}

func main() {
	lambda.Start(handleRequest)
}
