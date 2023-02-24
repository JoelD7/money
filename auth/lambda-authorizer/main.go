package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/JoelD7/money/api/shared/env"
	"github.com/JoelD7/money/api/shared/restclient"
	"github.com/JoelD7/money/api/shared/secrets"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt/v3"
	"log"
	"math/big"
	"os"
	"strings"
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
	errInvalidToken       = errors.New("invalid token")
	errSigningKeyNotFound = errors.New("signing key not found")
	errUnauthorized       = errors.New("Unauthorized")

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

var (
	kidSecretName = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	awsRegion     = env.GetString("REGION", "us-east-1")
	jwtAudience   = env.GetString("AUDIENCE", "https://localhost:3000")
	jwtIssuer     = env.GetString("ISSUER", "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging")
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

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := strings.ReplaceAll(event.AuthorizationToken, "Bearer ", "")

	payload, err := getTokenPayload(token)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errUnauthorized
	}

	err = verifyToken(payload, token)
	if err != nil {
		return defaultDenyAllPolicy(event.MethodArn), err
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

func getTokenPayload(token string) (*jwtPayload, error) {
	var payload *jwtPayload

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 3 {
		errorLogger.Println(errInvalidToken)
		return nil, errInvalidToken
	}

	payloadPart, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		errorLogger.Println(err)
		return nil, err
	}

	err = json.Unmarshal(payloadPart, &payload)
	if err != nil {
		errorLogger.Println(err)
		return nil, err
	}

	return payload, nil
}

func defaultDenyAllPolicy(methodArn string) events.APIGatewayCustomAuthorizerResponse {
	tmp := strings.Split(methodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]

	resp := NewAuthorizerResponse("user", awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]
	resp.DenyAllMethods()

	return resp.APIGatewayCustomAuthorizerResponse
}

func verifyToken(payload *jwtPayload, token string) error {
	response, err := restclient.Get(payload.Issuer + "/auth/jwks")
	if err != nil {
		return err
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			errorLogger.Println(closeErr)
			err = closeErr
		}
	}()

	jwksVal := new(jwks)
	err = json.NewDecoder(response.Body).Decode(jwksVal)
	if err != nil {
		errorLogger.Println(err)
		return err
	}

	publicKey, err := getPublicKey(jwksVal)
	if err != nil {
		errorLogger.Println(err)
		return err
	}

	decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(publicKey))
	receivedPayload := &jwt.Payload{}

	err = validateJWTPayload(token, receivedPayload, decryptingHash)

	return err
}

func validateJWTPayload(token string, payload *jwt.Payload, decryptingHash *jwt.RSASHA) error {
	//now := time.Now()

	//nbfValidator := jwt.NotBeforeValidator(now)
	//expValidator := jwt.ExpirationTimeValidator(now)
	issValidator := jwt.IssuerValidator(jwtIssuer)
	audValidator := jwt.AudienceValidator(jwt.Audience{jwtAudience})

	validatePayload := jwt.ValidatePayload(payload, issValidator, audValidator)

	_, err := jwt.Verify([]byte(token), decryptingHash, payload, validatePayload)
	if err != nil {
		errorLogger.Println(err)
		return err
	}

	return nil
}

func getPublicKey(jwksVal *jwks) (*rsa.PublicKey, error) {
	kid, err := getKidFromSecret()
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
		errorLogger.Println(errSigningKeyNotFound)
		return nil, errSigningKeyNotFound
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(signingKey.N)
	if err != nil {
		errorLogger.Println(err)
		return nil, err
	}

	n := new(big.Int)
	n.SetBytes(nBytes)

	eBytes, err := base64.RawURLEncoding.DecodeString(signingKey.E)
	if err != nil {
		errorLogger.Println(err)
		return nil, err
	}

	e := new(big.Int)
	e.SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func getKidFromSecret() (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
		errorLogger.Println(err)
		return "", err
	}

	return *kidSecret.SecretString, nil

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
