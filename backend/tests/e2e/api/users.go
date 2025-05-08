package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/JoelD7/money/backend/models"
)

// GetMe returns data about the requesting user
func (e *E2ERequester) GetMe(t *testing.T) (*models.User, error) {
	url := fmt.Sprintf("%s/users", e.baseUrl)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("user request building failed: %w", err)
	}

	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("user request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(res.StatusCode, res.Body)
	}

	var user models.User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user response decoding failed: %w", err)
	}

	return &user, nil
}

func (e *E2ERequester) DeleteUser(username string, t *testing.T) (int, error) {
	url := fmt.Sprintf("%s/users/%s", e.baseUrl, username)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return 0, fmt.Errorf("user deletion request building failed: %w", err)
	}

	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("user deletion request failed: %w", err)
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

func (e *E2ERequester) CreateCategory(category *models.Category, headers map[string]string, t *testing.T) (int, error) {
	var createdCategory models.Category

	t.Cleanup(func() {
		if createdCategory.ID == "" {
			return
		}

		statusCode, err := e.DeleteCategory(createdCategory.ID, t)
		if statusCode != http.StatusNoContent || err != nil {
			t.Logf("Failed to delete category %s: %v", createdCategory.ID, err)
		}
	})

	url := fmt.Sprintf("%s/users/categories", e.baseUrl)

	requestBody, err := json.Marshal(category)
	if err != nil {
		return 0, fmt.Errorf("category creation request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return 0, fmt.Errorf("category creation request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("category creation request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	err = json.NewDecoder(res.Body).Decode(&createdCategory)
	if err != nil {
		return res.StatusCode, fmt.Errorf("category creation response decoding failed: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	return res.StatusCode, nil
}

func (e *E2ERequester) DeleteCategory(categoryID string, t *testing.T) (int, error) {
	url := fmt.Sprintf("%s/users/categories/%s", e.baseUrl, categoryID)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return 0, fmt.Errorf("category deletion request building failed: %w", err)
	}

	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("category deletion request failed: %w", err)
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
