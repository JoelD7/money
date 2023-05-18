package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/hash"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/restclient"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/storage/invalidtoken"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt/v3"
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
	Stage string
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
	errInvalidToken       = errors.New("invalid_access_token")
	errSigningKeyNotFound = errors.New("signing key not found")
	errUnauthorized       = errors.New("Unauthorized")
)

var (
	kidSecretName = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	awsRegion     = env.GetString("REGION", "us-east-1")
	jwtAudience   = env.GetString("AUDIENCE", "https://localhost:3000")
	jwtIssuer     = env.GetString("ISSUER", "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging")

	invalidJWTErrs = []error{jwt.ErrAudValidation, jwt.ErrExpValidation, jwt.ErrIatValidation, jwt.ErrIssValidation,
		jwt.ErrJtiValidation, jwt.ErrNbfValidation, jwt.ErrSubValidation}
)

type jwtPayload struct {
	Scope string `json:"scope"`
	*jwt.Payload
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid,omitempty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type request struct {
	log          logger.LogAPI
	startingTime time.Time
}

func (req *request) init() {
	req.startingTime = time.Now()
}

func (req *request) finish() {
	req.log.LogLambdaTime(req.startingTime, recover())
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

	payload, err := req.getTokenPayload(token)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errUnauthorized
	}

	err = req.verifyToken(ctx, payload, token)
	if err != nil {
		return defaultDenyAllPolicy(event.MethodArn, err), nil
	}

	principalID := payload.Subject

	tmp := strings.Split(event.MethodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]

	resp := NewAuthorizerResponse(principalID, awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]
	resp.AllowAllMethods()

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func (req *request) getTokenPayload(token string) (*jwtPayload, error) {
	var payload *jwtPayload

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 3 {
		req.log.Error("invalid_token_length_detected", errInvalidToken, []logger.Object{})

		return nil, errInvalidToken
	}

	payloadPart, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		req.log.Error("payload_decoding_failed", err, []logger.Object{})

		return nil, err
	}

	err = json.Unmarshal(payloadPart, &payload)
	if err != nil {
		req.log.Error("payload_unmarshalling_failed", err, []logger.Object{})

		return nil, err
	}

	return payload, nil
}

func defaultDenyAllPolicy(methodArn string, err error) events.APIGatewayCustomAuthorizerResponse {
	tmp := strings.Split(methodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]

	resp := NewAuthorizerResponse("user", awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]
	resp.DenyAllMethods()

	resp.APIGatewayCustomAuthorizerResponse.Context = map[string]interface{}{
		"stringKey": err.Error(),
	}

	return resp.APIGatewayCustomAuthorizerResponse
}

func (req *request) verifyToken(ctx context.Context, payload *jwtPayload, token string) error {
	response, err := restclient.Get(payload.Issuer + "/auth/jwks")
	if err != nil {
		req.log.Error("getting_jwks_failed", err, []logger.Object{})
		return err
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			req.log.Error("closing_response_body_failed", closeErr, []logger.Object{})

			err = closeErr
		}
	}()

	jwksVal := new(jwks)
	err = json.NewDecoder(response.Body).Decode(jwksVal)
	if err != nil {
		req.log.Error("decoding_response_body_failed", err, []logger.Object{})

		return err
	}

	publicKey, err := req.getPublicKey(jwksVal)
	if err != nil {
		req.log.Error("getting_public_key_failed", err, []logger.Object{})

		return err
	}

	decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(publicKey))
	receivedPayload := &jwt.Payload{}

	err = req.validateJWTPayload(token, receivedPayload, decryptingHash)
	if err != nil {
		return err
	}

	err = req.compareAccessTokenAgainstBlacklist(ctx, payload.Subject, token)
	if errors.Is(err, errInvalidToken) {
		req.log.Warning("blacklisted_token_use_detected", err, []logger.Object{
			logger.MapToLoggerObject("token", map[string]interface{}{
				"s_value": token,
			}),
		})
	}

	return err
}

func (req *request) validateJWTPayload(token string, payload *jwt.Payload, decryptingHash *jwt.RSASHA) error {
	now := time.Now()

	expValidator := jwt.ExpirationTimeValidator(now)
	issValidator := jwt.IssuerValidator(jwtIssuer)
	audValidator := jwt.AudienceValidator(jwt.Audience{jwtAudience})

	validatePayload := jwt.ValidatePayload(payload, issValidator, audValidator, expValidator)

	_, err := jwt.Verify([]byte(token), decryptingHash, payload, validatePayload)
	if isErrorInvalidJWT(err) {
		req.log.Error("invalid_jwt", err, []logger.Object{
			logger.MapToLoggerObject("jwt_payload", map[string]interface{}{
				"s_subject":    payload.Subject,
				"s_audience":   payload.Audience,
				"f_expiration": payload.ExpirationTime,
			}),
		})

		return errInvalidToken
	}

	if err != nil {
		req.log.Error("jwt_validation_failed", err, []logger.Object{})

		return err
	}

	return nil
}

func isErrorInvalidJWT(err error) bool {
	for _, e := range invalidJWTErrs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func (req *request) compareAccessTokenAgainstBlacklist(ctx context.Context, email, token string) error {
	invalidTokens, err := invalidtoken.GetAllForPerson(ctx, email)
	if errors.Is(err, invalidtoken.ErrNotFound) {
		req.log.Info("no_tokens_found_for_user", []logger.Object{
			logger.MapToLoggerObject("person_data", map[string]interface{}{
				"s_email": email,
			}),
		})

		return nil
	}

	if err != nil {
		return err
	}

	for _, it := range invalidTokens {
		if it.Type == string(invalidtoken.TypeRefresh) {
			continue
		}

		err = hash.CompareWithToken(it.Token, token)
		if err == nil {
			req.log.Warning("invalid_token_use_detected", errInvalidToken, []logger.Object{
				logger.MapToLoggerObject("token_comparison", map[string]interface{}{
					"s_bearer_token":             token,
					"s_saved_invalid_token_hash": it.Token,
				}),
			})

			return errInvalidToken
		}

		if !errors.Is(err, hash.ErrHashMismatch) {
			req.log.Error("token_comparison_against_blacklist_failed", err, []logger.Object{
				logger.MapToLoggerObject("token_comparison", map[string]interface{}{
					"s_bearer_token":             token,
					"s_saved_invalid_token_hash": it.Token,
				}),
			})
		}
	}

	return nil
}

func (req *request) getPublicKey(jwksVal *jwks) (*rsa.PublicKey, error) {
	kid, err := req.getKidFromSecret()
	if err != nil {
		return nil, err
	}

	var signingKey *jwk

	for _, key := range jwksVal.Keys {
		if key.Kid == kid {
			signingKey = &key
		}
	}

	if signingKey == nil {
		req.log.Error("signing_key_not_found", errSigningKeyNotFound, []logger.Object{
			logger.MapToLoggerObject("kid", map[string]interface{}{
				"s_secret": kid,
			}),
		})

		return nil, errSigningKeyNotFound
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(signingKey.N)
	if err != nil {
		return nil, err
	}

	n := new(big.Int)
	n.SetBytes(nBytes)

	eBytes, err := base64.RawURLEncoding.DecodeString(signingKey.E)
	if err != nil {
		return nil, err
	}

	e := new(big.Int)
	e.SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func (req *request) getKidFromSecret() (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
		return "", err
	}

	return kidSecret, nil

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

func NewAuthorizerResponse(principalID string, AccountID string) *AuthorizerResponse {
	return &AuthorizerResponse{
		APIGatewayCustomAuthorizerResponse: events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: principalID,
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
			},
		},
		Region:    awsRegion,
		AccountID: AccountID,

		// Replace the placeholder value with a default API Gateway API id to be used in the policy.
		// Beware of using '*' since it will not simply mean any API Gateway API id, because stars will greedily expand over '/' or other separators.
		// See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_resource.html for more details.
		APIID: "<<restApiId>>",

		// Replace the placeholder value with a default stage to be used in the policy.
		// Beware of using '*' since it will not simply mean any stage, because stars will greedily expand over '/' or other separators.
		// See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_resource.html for more details.
		Stage: "<<stage>>",
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
