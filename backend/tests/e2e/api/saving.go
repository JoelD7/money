package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"testing"
)

type savingsResponse struct {
	Savings []*models.Saving `json:"savings"`
	NextKey string           `json:"next_key"`
}

func (e *E2ERequester) CreateSaving(saving *models.Saving, t *testing.T) (*models.Saving, int, error) {
	var createdSaving models.Saving

	t.Cleanup(func() {
		if createdSaving.SavingID == "" {
			return
		}

		statusCode, err := e.DeleteSaving(createdSaving.SavingID)
		if statusCode != http.StatusNoContent || err != nil {
			t.Logf("Failed to delete saving %s: %v", createdSaving.SavingID, err)
		}
	})

	requestBody, err := json.Marshal(saving)
	if err != nil {
		return nil, 0, fmt.Errorf("saving request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+savingsEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("saving request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", uuid.NewString())
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("saving request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusCreated {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	err = json.NewDecoder(res.Body).Decode(&createdSaving)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("saving response decoding failed: %w", err)
	}

	return &createdSaving, res.StatusCode, nil
}

func (e *E2ERequester) GetSavings(params *models.QueryParameters) ([]*models.Saving, string, int, error) {
	request, err := http.NewRequest(http.MethodGet, e.baseUrl+savingsEndpoint, nil)
	if err != nil {
		return nil, "", 0, fmt.Errorf("savings request building failed: %w", err)
	}

	q := request.URL.Query()
	params.ParseAsURLValues(&q)

	request.URL.RawQuery = q.Encode()
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, "", 0, fmt.Errorf("savings request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, "", res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	var response savingsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, "", res.StatusCode, fmt.Errorf("savings response decoding failed: %w", err)
	}

	return response.Savings, response.NextKey, res.StatusCode, nil
}

func (e *E2ERequester) DeleteSaving(savingID string) (int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, savingsEndpoint, savingID)
	if err != nil {
		return 0, fmt.Errorf("saving deletion endpoint building failed: %w", err)
	}

	return e.DeleteResource(endpoint)
}
