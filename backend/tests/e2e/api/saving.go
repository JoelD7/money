package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"net/http"
)

func (e *E2ERequester) CreateSaving(saving *models.Saving) (*models.Saving, int, error) {
	requestBody, err := json.Marshal(saving)
	if err != nil {
		return nil, 0, fmt.Errorf("saving request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+savingsEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("saving request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
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

	var createdSaving models.Saving
	err = json.NewDecoder(res.Body).Decode(&createdSaving)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("saving response decoding failed: %w", err)
	}

	return &createdSaving, res.StatusCode, nil
}

func (e *E2ERequester) DeleteSaving(savingID string) (int, error) {
	request, err := http.NewRequest(http.MethodDelete, e.baseUrl+savingsEndpoint+"/"+savingID, nil)
	if err != nil {
		return 0, fmt.Errorf("saving deletion request building failed: %w", err)
	}

	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("saving deletion request failed: %w", err)
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
