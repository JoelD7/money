package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

func (e *E2ERequester) CreatePeriod(period *models.Period, t TestCleaner) (*models.Period, int, error) {
	var createdPeriod models.Period

	defer t.Cleanup(func() {
		//if createdPeriod.ID == "" {
		//	return
		//}
		//
		//status, err := e.DeletePeriod(createdPeriod.ID)
		//if status != http.StatusNoContent || err != nil {
		//	t.Logf("Failed to delete period %s: %v", createdPeriod.ID, err)
		//}
	})

	requestBody, err := json.Marshal(period)
	if err != nil {
		return nil, 0, fmt.Errorf("period request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+periodsEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("period request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", uuid.NewString())
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

	err = json.NewDecoder(res.Body).Decode(&createdPeriod)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("period response decoding failed: %w", err)
	}

	return &createdPeriod, res.StatusCode, nil
}

func (e *E2ERequester) DeletePeriod(periodID string) (int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, periodsEndpoint, periodID)
	if err != nil {
		return 0, fmt.Errorf("period deletion endpoint building failed: %w", err)
	}

	return e.DeleteResource(endpoint)
}
