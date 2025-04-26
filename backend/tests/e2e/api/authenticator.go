package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func (e *E2ERequester) SignUp(username, fullname, password string, headers map[string]string, t *testing.T) (statusCode int, err error) {
	isUserCreated := false

	t.Cleanup(func() {
		if !isUserCreated {
			return
		}

		deleteStatusCode, deleteErr := e.DeleteUser(username, t)
		if deleteErr != nil {
			statusCode = deleteStatusCode
			err = fmt.Errorf("couldn't delete user '%s' during cleanup: %w", username, deleteErr)
		}
	})

	requestBody := []byte(fmt.Sprintf(`{"username":"%s","fullname":"%s","password":"%s"}`, username, fullname, password))

	url := fmt.Sprintf("%s/auth/signup", e.baseUrl)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return 0, fmt.Errorf("sign up request building failed: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := e.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("sign up request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusCreated {
		return res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	isUserCreated = true

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("sign up response body reading failed: %w", err)
	}

	var signupResponse authResponse
	err = json.Unmarshal(responseBody, &signupResponse)
	if err != nil {
		return 0, fmt.Errorf("signup response unmarshalling failed: %w", err)
	}

	e.accessToken = signupResponse.AccessToken

	return res.StatusCode, nil
}
