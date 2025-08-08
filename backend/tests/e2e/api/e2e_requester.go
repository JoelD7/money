package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"io"
	"net/http"
)

var (
	loginEndpoint       = "/auth/login"
	savingsEndpoint     = "/savings"
	expensesEndpoint    = "/expenses"
	usersEndpoint       = "/users"
	savingGoalsEndpoint = "/savings/goals"
	periodsEndpoint     = "/periods"
)

type TestCleaner interface {
	Cleanup(f func())
	Logf(format string, args ...any)
}

// E2ERequester is a type that will be used to make requests to the backend. It's main purpose is to hold the access token.
type E2ERequester struct {
	accessToken string
	baseUrl     string
	client      *http.Client

	Username string
}

type ErrorResponse struct {
	Code     int    `json:"code"`
	HTTPCode int    `json:"http_code"`
	Message  string `json:"message"`
}

type authResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewE2ERequester() (*E2ERequester, error) {
	requester := &E2ERequester{
		client:   &http.Client{},
		baseUrl:  env.GetString("BASE_URL", ""),
		Username: env.GetString("E2E_USER_USERNAME", "e2e_test@mail.com"),
	}

	err := requester.login()
	if err != nil {
		return nil, err
	}

	return requester, nil
}

func (e *E2ERequester) login() error {
	//This user already exists in the DB. Was created with the sole purpose of using it for e2e tests.
	password := env.GetString("E2E_USER_PASSWORD", "")

	loginRequestBody := fmt.Sprintf(`{"username":"%s","password":"%s"}`, e.Username, password)

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+loginEndpoint, bytes.NewReader([]byte(loginRequestBody)))
	if err != nil {
		return fmt.Errorf("login request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := e.client.Do(request)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Printf("login response body closing failed: %v", err)
		}
	}()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("login response body reading failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("login request failed with status: %d", res.StatusCode)
	}

	var loginRes authResponse
	err = json.Unmarshal(responseBody, &loginRes)
	if err != nil {
		return fmt.Errorf("login response unmarshalling failed: %w", err)
	}

	e.accessToken = loginRes.AccessToken

	return nil
}

func (e *E2ERequester) DeleteResource(endpoint string) (int, error) {
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusNoContent {
		return res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	return res.StatusCode, nil
}
