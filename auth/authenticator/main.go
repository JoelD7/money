// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/JoelD7/money/api/shared/router"
	storage "github.com/JoelD7/money/api/storage/person"
	"github.com/JoelD7/money/auth/authenticator/secrets"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	errMissingEmail      = errors.New("missing email")
	errMissingPassword   = errors.New("missing password")
	errWrongCredentials  = errors.New("the email or password are incorrect")
	errInvalidEmail      = errors.New("email is invalid")
	errGettingKeysForJWT = errors.New("error getting keys for JWT")

	errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
)

const (
	passwordCost = bcrypt.DefaultCost
	emailRegex   = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"
	secretName   = "staging/money/authenticator/rsakeys/v3"
)

type signUpBody struct {
	FullName string `json:"fullname"`
	*Credentials
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type rsaSecret struct {
	Public  string `json:"public"`
	Private string `json:"private"`
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

	key, err := getRSAPrivateKey()
	if err != nil {
		return serverError(err)
	}

	fmt.Println("Private key valid? ", key.Validate())

	err = json.Unmarshal([]byte(request.Body), reqBody)
	if err != nil {
		return serverError(err)
	}

	err = validateCredentials(reqBody)
	if err != nil {
		return clientError(err)
	}

	person, err := storage.GetPerson(reqBody.Email)
	if err != nil {
		return clientError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(reqBody.Password))
	if err != nil {
		return clientError(errWrongCredentials)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

//func generateJWT() error {
//
//	now := time.Now()
//	var hs = jwt.NewRS256(jwt.RSAPublicKey(publicKey))
//	pl := jwt.Payload{
//		Issuer:         "gbrlsnchs",
//		Subject:        "someone",
//		Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
//		ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
//		NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
//		IssuedAt:       jwt.NumericDate(now),
//		JWTID:          "foobar",
//	}
//
//	token, err := jwt.Sign(pl, hs)
//	if err != nil {
//		// ...
//	}
//}

func getRSAPrivateKey() (*rsa.PrivateKey, error) {
	secret, err := secrets.GetSecret(context.Background(), secretName)
	if err != nil && errors.Is(err, secrets.ErrSecretNotFound) {
		fmt.Println("secret not found")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var rsa rsaSecret

	fmt.Println("keys: ", *secret.SecretString)

	err = json.Unmarshal([]byte(*secret.SecretString), &rsa)
	if err != nil {
		return nil, errGettingKeysForJWT
	}

	//PEM encoded private key
	privatePEMData := []byte(rsa.Private)
	fmt.Println(privatePEMData)
	//Create pem block
	privatePemBlock, _ := pem.Decode(privatePEMData)
	fmt.Println("privatePemBlock: ", privatePemBlock)
	if privatePemBlock == nil || !strings.Contains(privatePemBlock.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM privatePemBlock containing private key")
	}

	//Use bytes from pem block to get an rsa.PrivateKey
	privateKey, err := x509.ParsePKCS1PrivateKey(privatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

//func getSecretRSAKeys() *secretsmanager.GetSecretValueOutput {
//	secret, err := secrets.GetSecret(context.Background(), secretName)
//	if err != nil && errors.Is(err, secrets.ErrSecretNotFound) {
//		fmt.Println(err)
//
//	}
//}

func createRSAKeys(ctx context.Context) error {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publickey := &privatekey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publickey)

	rsa := &rsaSecret{
		Public:  string(publicKeyBytes),
		Private: string(privateKeyBytes),
	}

	rsaBytes, err := json.Marshal(rsa)
	if err != nil {
		return err
	}

	secrets.CreateSecret(ctx, secretName, "RSA Keys to sign and verify JWTs", rsaBytes)
	return nil
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
	})

	lambda.Start(route.Handle)
}
