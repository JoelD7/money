package authenticator

import (
	"net/http"
	"os"
	"testing"

	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(err)
	}

	logger.InitLogger(logger.ConsoleImplementation)

	os.Exit(m.Run())
}

func TestSignUp(t *testing.T) {
	c := require.New(t)

	e2eRequester, err := api.NewE2ERequester()
	c.NoError(err)
	c.NotNil(e2eRequester)

	username := "signup_test@mail.com"

	t.Run("Success", func(t *testing.T) {
		headers := map[string]string{
			"Idempotency-Key": "1234",
		}

		statusCode, err := e2eRequester.SignUp(username, "John Doe", "password", headers, t)
		c.NoError(err)
		c.Equal(http.StatusCreated, statusCode)

		user, err := e2eRequester.GetMe(t)
		c.NoError(err)
		c.Equal(username, user.Username)
	})

	t.Run("Missing idempotency key", func(t *testing.T) {
		headers := map[string]string{}

		statusCode, err := e2eRequester.SignUp(username, "John Doe", "password", headers, t)
		c.Error(err)
		c.Equal(http.StatusBadRequest, statusCode)
	})
}
