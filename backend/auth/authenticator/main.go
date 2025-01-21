// This lambda represents the authentication server.
// Authenticates users and generates JWTs.
package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/router"
	"github.com/JoelD7/money/backend/shared/uuid"
	"github.com/JoelD7/money/backend/shared/validate"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"strings"
)

var (
	errCookiesNotFound              = apigateway.NewError("cookies not found in request object", http.StatusBadRequest)
	errMissingRefreshTokenInCookies = apigateway.NewError("missing refresh token in cookies", http.StatusBadRequest)
	errUserNotFound                 = apigateway.NewError("", http.StatusBadRequest)
)

const (
	refreshTokenCookieName = "RefreshToken"
)

type signUpBody struct {
	FullName string `json:"fullname"`
	*Credentials
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (c *Credentials) Key() string {
	return "user_data"
}

func (c *Credentials) Value() map[string]interface{} {
	return map[string]interface{}{
		"username": c.Username,
	}
}

func getRefreshTokenCookie(request *apigateway.Request) (string, error) {
	cookies, ok := request.Headers["Cookie"]
	if !ok {
		return "", errCookiesNotFound
	}

	cookieParts := make([]string, 0)
	var name, value string

	for _, cookie := range strings.Split(cookies, ";") {
		cookieParts = strings.Split(cookie, "=")
		if len(cookieParts) < 2 {
			continue
		}

		name = strings.TrimSpace(cookieParts[0])
		value = strings.TrimSpace(cookieParts[1])

		if name == "" || value == "" {
			continue
		}

		if name == refreshTokenCookieName && len(cookieParts) > 1 {
			return value, nil
		}
	}

	return "", errMissingRefreshTokenInCookies
}

func validateCredentials(email, password string) error {
	err := validate.Email(email)
	if err != nil {
		return err
	}

	if password == "" {
		return models.ErrMissingPassword
	}

	return nil
}

func main() {
	envConfig, err := env.LoadEnv(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to load environment variables: %w", err))
	}

	rootRouter := router.NewRouter(envConfig)

	rootRouter.Route("/auth", func(r *router.Router) {
		r.Post("/login", logInHandler)

		r.Post("/signup", signUpHandler)

		r.Post("/token", tokenHandler)
		r.Get("/jwks", jwksHandler)
		r.Post("/logout", logoutHandler)
	})

	lambda.Start(func(ctx context.Context, request *apigateway.Request) (res *apigateway.Response, err error) {
		logger.InitLogger(logger.LogstashImplementation)
		logger.AddToContext("request_id", uuid.Generate(request.RequestContext.ExtendedRequestID))

		defer func() {
			err = logger.Finish()
			if err != nil {
				panic(fmt.Errorf("failed to finish logger: %w", err))
			}
		}()

		return rootRouter.Handle(ctx, request)
	})
}
