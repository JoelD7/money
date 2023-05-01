// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/JoelD7/money/backend/entities"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/shared/utils"
	storagePerson "github.com/JoelD7/money/backend/storage/person"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt/v3"
)

var (
	errMissingEmail     = errors.New("missing email")
	errMissingPassword  = errors.New("missing password")
	errWrongCredentials = errors.New("the email or password are incorrect")
	errInvalidEmail     = errors.New("email is invalid")
	errCookiesNotFound  = errors.New("cookies not found in request object")
)

var (
	accessTokenAudience  = env.GetString("TOKEN_AUDIENCE", "https://localhost:3000")
	accessTokenIssuer    = env.GetString("TOKEN_ISSUER", "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging")
	accessTokenScope     = env.GetString("TOKEN_SCOPE", "read write")
	privateSecretName    = env.GetString("TOKEN_PRIVATE_SECRET", "staging/money/rsa/private")
	publicSecretName     = env.GetString("TOKEN_PUBLIC_SECRET", "staging/money/rsa/public")
	kidSecretName        = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	accessTokenDuration  = env.GetInt("ACCESS_TOKEN_DURATION", 300)
	refreshTokenduration = env.GetInt("REFRESH_TOKEN_DURATION", 2592000)

	redisEndpoint = env.GetString("REDIS_ENDPOINT", "money-redis-cluster-ro.sfenn7.ng.0001.use1.cache.amazonaws.com:6379")
)

const (
	passwordCost = bcrypt.DefaultCost
	emailRegex   = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"
)

type signUpBody struct {
	FullName string `json:"full_name"`
	*Credentials
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
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

type jwtPayload struct {
	Scope string `json:"scope,omitempty"`
	*jwt.Payload
}

type requestHandler struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	log          logger.LogAPI
	startingTime time.Time
}

func (c *Credentials) LogName() string {
	return "user_data"
}

func (c *Credentials) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"email": c.Email,
	}
}

func (req *requestHandler) init() {
	req.startingTime = time.Now()
}

func (req *requestHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, recover())
}

func refreshTokenHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req := &requestHandler{
		log: logger.NewLoggerWithHandler("refresh-token"),
	}

	req.init()
	defer req.finish()

	return req.processRefreshToken(request)
}

func (req *requestHandler) processRefreshToken(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var credentials Credentials
	ctx := context.Background()

	err := json.Unmarshal([]byte(request.Body), &credentials)
	if err != nil {
		req.log.Error("request_body_unmarshal_failed", err, []logger.Object{})

		return req.serverError()
	}

	cookies, ok := request.Headers["Cookie"]
	if !ok {
		req.log.Error("cookies_not_found_in_request", errCookiesNotFound, []logger.Object{})

		return req.serverError()
	}

	var refreshToken string

	for _, cookie := range strings.Split(cookies, ";") {
		if strings.HasPrefix("refresh_token", cookie) {
			refreshToken = strings.Split(cookie, "=")[1]
			break
		}
	}

	person, err := storagePerson.GetPersonByEmail(ctx, credentials.Email)
	if err != nil {
		req.log.Error("fetching_user_from_storage_failed", err, []logger.Object{})

		return req.serverError()
	}

	if person.PreviousRefreshToken != "" && refreshToken == person.PreviousRefreshToken {
		//	TODO: invalidate all tokens of the user, forcing him to reauthenticate

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Invalid credentials",
		}, nil
	}

	if refreshToken == person.RefreshToken {
		// TODO: issue new access token and refresh token
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}

func signUpHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req := &requestHandler{
		log: logger.NewLoggerWithHandler("sign-up"),
	}

	req.init()
	defer req.finish()

	return req.processSignUp(request)
}

func (req *requestHandler) processSignUp(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	reqBody := &signUpBody{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		req.log.Error("request_body_json_unmarshal_failed", err, []logger.Object{})

		return req.serverError()
	}

	err = req.validateCredentials(reqBody.Credentials)
	if err != nil {
		req.log.Error("credentials_validation_failed", err, []logger.Object{})

		return req.clientError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), passwordCost)
	if err != nil {
		req.log.Error("password_hashing_failed", err, []logger.Object{})

		return req.serverError()
	}

	err = storagePerson.CreatePerson(ctx, reqBody.FullName, reqBody.Email, string(hashedPassword))
	if err != nil && errors.Is(err, storagePerson.ErrExistingUser) {
		req.log.Warning("user_creation_failed", err, []logger.Object{})

		return req.clientError(err)
	}

	if err != nil {
		req.log.Error("sign_up_process_failed", err, []logger.Object{})

		return req.serverError()
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func logInHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req := &requestHandler{
		log: logger.NewLoggerWithHandler("log-in"),
	}

	req.init()
	defer req.finish()

	return req.processLogin(request)
}

func (req *requestHandler) processLogin(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx := context.Background()

	reqBody := &Credentials{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		req.log.Error("request_body_json_unmarshal_failed", err, []logger.Object{})

		return req.serverError()
	}

	err = req.validateCredentials(reqBody)
	if err != nil {
		req.log.Error("credentials_validation_failed", err, []logger.Object{})

		return req.clientError(err)
	}

	person, err := storagePerson.GetPersonByEmail(ctx, reqBody.Email)
	if err != nil {
		req.log.Error("user_fetching_failed", err, []logger.Object{})

		return req.clientError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(reqBody.Password))
	if err != nil {
		req.log.Error("password_mismatch", err, []logger.Object{reqBody})

		return req.clientError(errWrongCredentials)
	}

	headers, err := req.setTokens(ctx, person)
	if err != nil {
		req.log.Error("token_setting_failed", err, []logger.Object{reqBody})

		return req.serverError()
	}

	responseBody, err := utils.GetJsonString(&loginResponse{req.AccessToken})
	if err != nil {
		req.log.Error("response_building_failed", err, []logger.Object{reqBody})

		return req.serverError()
	}

	req.log.Info("login_succeeded", []logger.Object{})

	return &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Body:              responseBody,
		MultiValueHeaders: headers,
	}, nil
}

func (req *requestHandler) setTokens(ctx context.Context, person *entities.Person) (map[string][]string, error) {
	now := time.Now()

	accessTokenExpiry := jwt.NumericDate(now.Add(time.Duration(accessTokenDuration) * time.Second))

	accessTokenPayload := &jwt.Payload{
		Issuer:         accessTokenIssuer,
		Subject:        person.Email,
		Audience:       jwt.Audience{accessTokenAudience},
		ExpirationTime: accessTokenExpiry,
		IssuedAt:       jwt.NumericDate(now),
	}

	accessToken, err := req.generateJWT(accessTokenPayload, accessTokenScope)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := jwt.NumericDate(now.Add(time.Duration(refreshTokenduration) * time.Second))

	refreshTokenPayload := &jwt.Payload{
		Subject:        person.Email,
		ExpirationTime: refreshTokenExpiry,
	}

	refreshToken, err := req.generateJWT(refreshTokenPayload, "")
	if err != nil {
		return nil, err
	}

	req.AccessToken = accessToken
	req.RefreshToken = refreshToken

	person.PreviousRefreshToken = person.RefreshToken
	person.RefreshToken = req.RefreshToken

	setCookieHeaders := map[string][]string{
		"Set-Cookie": {
			fmt.Sprintf("access_token=%s; Path=/; Expires=%s; Secure; HttpOnly", req.AccessToken,
				accessTokenExpiry.Format(time.RFC1123)),
			fmt.Sprintf("refresh_token=%s; Path=/; Expires=%s; Secure; HttpOnly", req.RefreshToken,
				refreshTokenExpiry.Format(time.RFC1123)),
		},
	}

	return setCookieHeaders, storagePerson.UpdatePerson(ctx, person)
}

func jwksHandler(_ *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req := &requestHandler{
		log: logger.NewLogger(),
	}

	req.init()
	defer req.finish()

	return req.processJWKS()
}

func (req *requestHandler) processJWKS() (*events.APIGatewayProxyResponse, error) {
	publicKey, err := req.getPublicKey()
	if err != nil {
		req.log.Error("public_key_fetching_failed", err, []logger.Object{})

		return req.serverError()
	}

	kid, err := req.getKidFromSecret()
	if err != nil {
		req.log.Error("kid_fetching_failed", err, []logger.Object{})

		return req.serverError()
	}

	response := jwks{
		[]jwk{
			{
				Kid: kid,
				Kty: "RSA",
				Use: "sig",
				N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
			},
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		req.log.Error("jwks_response_marshall_failed", err, []logger.Object{})

		return req.serverError()
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
	}, nil
}

func (req *requestHandler) generateJWT(payload *jwt.Payload, scope string) (string, error) {
	priv, err := req.getPrivateKey()
	if err != nil {
		req.log.Error("private_key_fetching_failed", err, []logger.Object{})

		return "", err
	}

	var signingHash = jwt.NewRS256(jwt.RSAPrivateKey(priv))

	p := jwtPayload{
		Scope:   scope,
		Payload: payload,
	}

	token, err := jwt.Sign(p, signingHash)
	if err != nil {
		req.log.Error("jwt_signing_failed", err, []logger.Object{})

		return "", err
	}

	//decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(pub))
	//receivedPayload := &jwt.Payload{}
	//
	//hd, err := jwt.Verify(Token, decryptingHash, receivedPayload)
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Println("Successfully verified. Header: ", hd)

	return string(token), nil
}

func (req *requestHandler) getPrivateKey() (*rsa.PrivateKey, error) {
	privateSecret, err := secrets.GetSecret(context.Background(), privateSecretName)
	if err != nil {
		return nil, err
	}

	privatePemBlock, _ := pem.Decode([]byte(privateSecret))
	if privatePemBlock == nil || !strings.Contains(privatePemBlock.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM private block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (req *requestHandler) getPublicKey() (*rsa.PublicKey, error) {
	publicSecret, err := secrets.GetSecret(context.Background(), publicSecretName)
	if err != nil {
		return nil, err
	}

	publicPemBlock, _ := pem.Decode([]byte(publicSecret))
	if publicPemBlock == nil || !strings.Contains(publicPemBlock.Type, "PUBLIC KEY") {
		return nil, fmt.Errorf("failed to decode PEM public block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// The kid of the JWK that contains the public key.
// Is stored in a secret so that the lambda-authorizer can have access to it to verify that the key received is the
// right one.
func (req *requestHandler) getKidFromSecret() (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
		return "", err
	}

	return kidSecret, nil

}

func (req *requestHandler) serverError() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func (req *requestHandler) clientError(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func (req *requestHandler) validateCredentials(login *Credentials) error {
	regex := regexp.MustCompile(emailRegex)

	if login.Email == "" {
		return errMissingEmail
	}

	if !regex.MatchString(login.Email) {
		return errInvalidEmail
	}

	if login.Password == "" {
		return errMissingPassword
	}

	return nil
}

func main() {
	route := router.NewRouter()

	route.Route("/auth", func(r *router.Router) {
		r.Post("/login", logInHandler)
		r.Post("/signup", signUpHandler)
		r.Post("/refresh-token", refreshTokenHandler)
		r.Get("/jwks", jwksHandler)
	})

	lambda.Start(route.Handle)
}
