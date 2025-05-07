package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/shared/uuid"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/shared"
	"github.com/JoelD7/money/backend/usecases"
	"strings"
	"sync"
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
	req  *requestInfo
	once sync.Once
)

type requestInfo struct {
	secretsManager secrets.SecretManager
	cacheRepo      cache.InvalidTokenManager
	startingTime   time.Time
	client         restclient.HttpClient
	err            error
}

func (req *requestInfo) init() {
	once.Do(func() {
		req.cacheRepo = cache.NewRedisCache()
		req.secretsManager = secrets.NewAWSSecretManager()
		req.client = restclient.New()
	})
	req.startingTime = time.Now()
	req.err = nil
}

func (req *requestInfo) finish() {
	defer func() {
		err := logger.Finish()
		if err != nil {
			logger.ErrPrintln("failed to finish logger", err)
		}
	}()

	logger.LogLambdaTime(req.startingTime, req.err, recover())
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (res events.APIGatewayCustomAuthorizerResponse, err error) {
	stackTrace, ctxError := shared.ExecuteLambda(ctx, func(ctx context.Context) {
		if req == nil {
			req = &requestInfo{}
		}

		req.init()
		defer req.finish()

		res, err = req.process(ctx, event)
	})

	if ctxError != nil {
		logger.Error("request_timeout", ctxError, req.getEventAsLoggerObject(event), models.Any("stack", map[string]interface{}{
			"s_trace": stackTrace,
		}))
	}

	return
}

func (req *requestInfo) process(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.ReplaceAll(event.AuthorizationToken, "Bearer ", "")

	verifyToken := usecases.NewTokenVerifier(req.client, req.secretsManager, req.cacheRepo)

	subject, err := verifyToken(ctx, token)
	if errors.Is(err, models.ErrUnauthorized) || errors.Is(err, models.ErrInvalidToken) {
		logger.Error("request_unauthorized", err, req.getEventAsLoggerObject(event))

		return events.APIGatewayCustomAuthorizerResponse{}, models.ErrUnauthorized
	}

	if err != nil {
		logger.Error("token_verification_failed", err, req.getEventAsLoggerObject(event))

		return events.APIGatewayCustomAuthorizerResponse{}, models.ErrUnauthorized
	}

	principalID := subject

	resp := NewAuthorizerResponse(event.MethodArn, principalID)

	resp.AllowAllMethods()

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func (req *requestInfo) getEventAsLoggerObject(event events.APIGatewayCustomAuthorizerRequest) models.LoggerField {
	return models.Any("authorizer_request", map[string]interface{}{
		"s_type":       event.Type,
		"s_method_arn": event.MethodArn,
	})
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
	awsRegion := env.GetString("AWS_REGION", "")

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
	_, err := env.LoadEnv(context.Background())
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	lambda.Start(func(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate(event.MethodArn))

		defer func() {
			err = logger.Finish()
			if err != nil {
				logger.ErrPrintln("failed to finish logger", err)
			}
		}()

		return handleRequest(ctx, event)
	})
}
