// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"context"
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
	errCookiesNotFound              = apigateway.NewError("cookies not found in request object", http.StatusBadRequest)
	errMissingRefreshTokenInCookies = apigateway.NewError("missing refresh token in cookies", http.StatusBadRequest)
	errUserNotFound                 = apigateway.NewError("", http.StatusBadRequest)
)

var (
	privateSecretName = env.GetString("TOKEN_PRIVATE_SECRET", "staging/money/rsa/private")
	publicSecretName  = env.GetString("TOKEN_PUBLIC_SECRET", "staging/money/rsa/public")
	kidSecretName     = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	awsRegion         = env.GetString("REGION", "us-east-1")
)

const (
	refreshTokenCookieName = "RefreshToken"
)

type signUpBody struct {
	FullName string `json:"full_name"`
	*Credentials
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (c *Credentials) LogName() string {
	return "user_data"
}

func (c *Credentials) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"username": c.Username,
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
