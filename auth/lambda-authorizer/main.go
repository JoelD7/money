// Copyright 2015-2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License").
// You may not use this file except in compliance with the License.
// A copy of the License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/shared/env"
	"github.com/JoelD7/money/auth/authenticator/secrets"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt/v3"
	"log"
	"math/big"
	"net/http"
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

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

var (
	kidSecretName = env.GetString("KID_SECRET", "staging/money/rsa/kid")
)

type jwtPayload struct {
	Scope          string       `json:"scope"`
	Issuer         string       `json:"iss,omitempty"`
	Subject        string       `json:"sub,omitempty"`
	Audience       jwt.Audience `json:"aud,omitempty"`
	ExpirationTime *jwt.Time    `json:"exp,omitempty"`
	NotBefore      *jwt.Time    `json:"nbf,omitempty"`
	IssuedAt       *jwt.Time    `json:"iat,omitempty"`
	JWTID          string       `json:"jti,omitempty"`
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
	//log.Println("Client token: " + event.AuthorizationToken)
	//log.Println("Method ARN: " + event.MethodArn)

	token := strings.Replace(event.AuthorizationToken, "Bearer", "", -1)
	payload, err := getTokenPayload(token)
	if err != nil {
		errorLogger.Println(err)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	err = verifyToken(payload, token)
	if err != nil {
		errorLogger.Println(err)

		return defaultDenyAllPolicy(event.MethodArn), err
	}

	principalID := payload.Subject

	// you can send a 401 Unauthorized response to the client by failing like so:
	// return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")

	// if the token is valid, a policy must be generated which will allow or deny access to the client

	// if access is denied, the client will recieve a 403 Access Denied response
	// if access is allowed, API Gateway will proceed with the backend integration configured on the method that was called

	// this function must generate a policy that is associated with the recognized principal user identifier.
	// depending on your use case, you might store policies in a DB, or generate them on the fly

	// keep in mind, the policy is cached for 5 minutes by default (TTL is configurable in the authorizer)
	// and will apply to subsequent calls to any method/resource in the RestApi
	// made with the same token

	//the example policy below denies access to all resources in the RestApi
	tmp := strings.Split(event.MethodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]

	resp := NewAuthorizerResponse(principalID, awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]
	resp.DenyAllMethods()
	// resp.AllowMethod(Get, "/pets/*")

	// new! -- add additional key-value pairs associated with the authenticated principal
	// these are made available by APIGW like so: $context.authorizer.<key>
	// additional context is cached
	resp.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func getTokenPayload(token string) (*jwtPayload, error) {
	var payload *jwtPayload

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 3 {
		return nil, errInvalidToken
	}

	payloadPart, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payloadPart, &payload)
	if err != nil {
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
	response, err := http.Get(payload.Issuer + "/auth/jwks")
	if err != nil {
		return err
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	var jwksVal jwks
	err = json.NewDecoder(response.Body).Decode(&jwksVal)
	if err != nil {
		return err
	}

	publicKey, err := getPublicKey(&jwksVal)
	if err != nil {
		return err
	}

	decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(publicKey))
	receivedPayload := &jwt.Payload{}

	hd, err := jwt.Verify([]byte(token), decryptingHash, receivedPayload)
	if err != nil {
		return err
	}

	fmt.Println("Successfully verified. Header: ", hd)
	fmt.Println("Received payload: ", receivedPayload)

	return err
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

func getKidFromSecret() (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
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
		// Replace the placeholder value with a default region to be used in the policy.
		// Beware of using '*' since it will not simply mean any region, because stars will greedily expand over '/' or other separators.
		// See https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_resource.html for more details.
		Region:    "<<region>>",
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
