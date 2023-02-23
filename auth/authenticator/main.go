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
	"github.com/JoelD7/money/api/shared/env"
	"math/big"

	"github.com/JoelD7/money/api/shared/router"
	storage "github.com/JoelD7/money/api/storage/person"
	"github.com/JoelD7/money/auth/authenticator/secrets"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gbrlsnchs/jwt/v3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	errMissingEmail      = errors.New("missing email")
	errMissingPassword   = errors.New("missing password")
	errWrongCredentials  = errors.New("the email or password are incorrect")
	errInvalidEmail      = errors.New("email is invalid")
	errGettingKeysForJWT = errors.New("error getting keys for JWT")

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

var (
	jwtAudience       = env.GetString("AUDIENCE", "https://localhost:3000")
	jwtIssuer         = env.GetString("ISSUER", "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging")
	jwtScope          = env.GetString("SCOPE", "read write")
	privateSecretName = env.GetString("PRIVATE_SECRET", "staging/money/rsa/private")
	publicSecretName  = env.GetString("PUBLIC_SECRET", "staging/money/rsa/public")
	kidSecretName     = env.GetString("KID_SECRET", "staging/money/rsa/kid")
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
	Password string `json:"password"`
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
	Scope string `json:"scope"`
	*jwt.Payload
}

func signUpHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBody := &signUpBody{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateCredentials(reqBody.Credentials)
	if err != nil {
		return clientError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), passwordCost)
	if err != nil {
		return serverError(err)
	}

	err = storage.CreatePerson(reqBody.FullName, reqBody.Email, string(hashedPassword))
	if err != nil && errors.Is(err, storage.ErrExistingUser) {
		return clientError(err)
	}

	if err != nil {
		return serverError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func logInHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBody := &Credentials{}

	err := json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateCredentials(reqBody)
	if err != nil {
		return clientError(err)
	}

	person, err := storage.GetPersonByEmail(reqBody.Email)
	if err != nil {
		return clientError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(reqBody.Password))
	if err != nil {
		return clientError(errWrongCredentials)
	}

	token, err := generateJWT(person.Email)
	if err != nil {
		return serverError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       token,
	}, nil
}

func jwksHandler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	publicKey, err := getPublicKey()
	if err != nil {
		return serverError(err)
	}

	kid, err := getKidFromSecret()
	if err != nil {
		return serverError(err)
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

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
	}, nil
}

func generateJWT(email string) (string, error) {
	now := time.Now()

	priv, err := getPrivateKey()
	if err != nil {
		return "", err
	}

	var signingHash = jwt.NewRS256(jwt.RSAPrivateKey(priv))

	payload := jwtPayload{
		Scope: jwtScope,
		Payload: &jwt.Payload{
			Issuer:   jwtIssuer,
			Subject:  email,
			Audience: jwt.Audience{jwtAudience},
			//ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
			//NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
			IssuedAt: jwt.NumericDate(now),
		},
	}

	token, err := jwt.Sign(payload, signingHash)
	if err != nil {
		return "", err
	}

	//decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(pub))
	//receivedPayload := &jwt.Payload{}
	//
	//hd, err := jwt.Verify(token, decryptingHash, receivedPayload)
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Println("Successfully verified. Header: ", hd)

	return string(token), nil
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	privateSecret, err := secrets.GetSecret(context.Background(), privateSecretName)
	if err != nil {
		return nil, err
	}

	privatePemBlock, _ := pem.Decode([]byte(*privateSecret.SecretString))
	fmt.Println(privatePemBlock)
	if privatePemBlock == nil || !strings.Contains(privatePemBlock.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM private block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKey() (*rsa.PublicKey, error) {
	publicSecret, err := secrets.GetSecret(context.Background(), publicSecretName)
	if err != nil {
		return nil, err
	}

	publicPemBlock, _ := pem.Decode([]byte(*publicSecret.SecretString))
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
func getKidFromSecret() (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
		return "", err
	}

	return *kidSecret.SecretString, nil

}

func serverError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(err error) (*events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func validateCredentials(login *Credentials) error {
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
		r.Get("/jwks", jwksHandler)
	})

	lambda.Start(route.Handle)
}
