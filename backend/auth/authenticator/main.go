// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
	"strings"
)

var (
	errCookiesNotFound              = errors.New("cookies not found in request object")
	errMissingRefreshTokenInCookies = errors.New("missing refresh token in cookies")

	responseByErrors = map[error]apigateway.Error{
		models.ErrMissingEmail:          {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingEmail.Error()},
		models.ErrInvalidEmail:          {HTTPCode: http.StatusBadRequest, Message: models.ErrInvalidEmail.Error()},
		models.ErrMissingPassword:       {HTTPCode: http.StatusBadRequest, Message: models.ErrMissingPassword.Error()},
		models.ErrInvalidToken:          {HTTPCode: http.StatusUnauthorized, Message: models.ErrInvalidToken.Error()},
		models.ErrMalformedToken:        {HTTPCode: http.StatusUnauthorized, Message: models.ErrMalformedToken.Error()},
		errMissingRefreshTokenInCookies: {HTTPCode: http.StatusBadRequest, Message: errMissingRefreshTokenInCookies.Error()},
		models.ErrExistingUser:          {HTTPCode: http.StatusBadRequest, Message: models.ErrExistingUser.Error()},
		models.ErrUserNotFound:          {HTTPCode: http.StatusBadRequest, Message: models.ErrUserNotFound.Error()},
		models.ErrWrongCredentials:      {HTTPCode: http.StatusBadRequest, Message: models.ErrWrongCredentials.Error()},
	}
)

var (
	privateSecretName = env.GetString("TOKEN_PRIVATE_SECRET", "staging/money/rsa/private")
	publicSecretName  = env.GetString("TOKEN_PUBLIC_SECRET", "staging/money/rsa/public")
	kidSecretName     = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	awsRegion         = env.GetString("REGION", "us-east-1")
)

const (
	accessTokenCookieName  = "AccessToken"
	refreshTokenCookieName = "RefreshToken"
)

type signUpBody struct {
	FullName string `json:"full_name"`
	*Credentials
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (c *Credentials) LogName() string {
	return "user_data"
}

func (c *Credentials) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"email": c.Email,
	}
}

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func getRefreshTokenCookie(request *apigateway.Request) (string, error) {
	cookies, ok := request.Headers["Cookie"]
	if !ok {
		return "", errCookiesNotFound
	}

	cookieParts := make([]string, 0)

	for _, cookie := range strings.Split(cookies, ";") {
		cookieParts = strings.Split(cookie, "=")

		if cookieParts[0] == "" || cookieParts[1] == "" {
			continue
		}

		if strings.HasPrefix(cookie, refreshTokenCookieName) && len(cookieParts) > 1 {
			return cookieParts[1], nil
		}
	}

	return "", errMissingRefreshTokenInCookies
}

func getErrorResponse(err error) (*apigateway.Response, error) {
	for mappedErr, responseErr := range responseByErrors {
		if errors.Is(err, mappedErr) {
			return apigateway.NewJSONResponse(responseErr.HTTPCode, responseErr.Message), nil
		}
	}

	return apigateway.NewErrorResponse(err), err
}

func main() {
	route := router.NewRouter()

	route.Route("/auth", func(r *router.Router) {
		r.Post("/login", logInHandler)
		r.Post("/signup", signUpHandler)
		r.Post("/token", tokenHandler)
		r.Get("/jwks", jwksHandler)
		r.Post("/logout", logoutHandler)
	})

	lambda.Start(route.Handle)
}
