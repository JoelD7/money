package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"net/http"
)

func (e *E2ERequester) CreatePeriod(period *models.Period) (*models.Period, int, error) {
	requestBody, err := json.Marshal(period)
	if err != nil {
		return nil, 0, fmt.Errorf("period request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+periodsEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("period request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("period request failed: %w", err)
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

	var createdPeriod models.Period
	err = json.NewDecoder(res.Body).Decode(&createdPeriod)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("period response decoding failed: %w", err)
	}

	return &createdPeriod, res.StatusCode, nil
}

func (e *E2ERequester) DeletePeriod(periodID string) (int, error) {
	request, err := http.NewRequest(http.MethodDelete, e.baseUrl+periodsEndpoint+"/"+periodID, nil)
	if err != nil {
		return 0, fmt.Errorf("period deletion request building failed: %w", err)
	}

	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("period deletion request failed: %w", err)
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
