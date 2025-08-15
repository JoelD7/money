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

const (
	incomesEndpoint = "/income"
)

// IncomesResponse represents the structure of the response when fetching multiple incomes.
type IncomesResponse struct {
	Incomes []*models.Income `json:"income"`
	NextKey string           `json:"next_key"`
}

// CreateIncome sends a request to create a new income.
// It also sets up a cleanup function to delete the income after the test.
func (e *E2ERequester) CreateIncome(income *models.Income, t TestCleaner) (*models.Income, int, error) {
	var createdIncome models.Income
	t.Cleanup(func() {
		if createdIncome.IncomeID == "" {
			return
		}

		statusCode, err := e.DeleteIncome(createdIncome.IncomeID)
		if statusCode != http.StatusNoContent || err != nil {
			t.Logf("Failed to delete income %s: %v", createdIncome.IncomeID, err)
		}
	})

	requestBody, err := json.Marshal(income)
	if err != nil {
		return nil, 0, fmt.Errorf("request body marshalling failed: %w", err)
	}

	endpoint, err := url.JoinPath(e.baseUrl, incomesEndpoint)
	if err != nil {
		return nil, 0, fmt.Errorf("request endpoint building failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", uuid.NewString())
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	err = json.NewDecoder(res.Body).Decode(&createdIncome)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("response decoding failed: %w", err)
	}

	return &createdIncome, res.StatusCode, nil
}

// GetIncomes sends a request to fetch multiple incomes based on the provided query parameters.
func (e *E2ERequester) GetIncomes(params *models.QueryParameters) ([]*models.Income, int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, incomesEndpoint)
	if err != nil {
		return nil, 0, fmt.Errorf("request endpoint building failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request building failed: %w", err)
	}

	addQueryParams(request, params, false)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	var incomesRes IncomesResponse
	err = json.NewDecoder(res.Body).Decode(&incomesRes)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("response decoding failed: %w", err)
	}

	return incomesRes.Incomes, res.StatusCode, nil
}

// DeleteIncome sends a request to delete a specific income by its ID.
func (e *E2ERequester) DeleteIncome(incomeID string) (int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, incomesEndpoint, incomeID)
	if err != nil {
		return 0, fmt.Errorf("income deletion endpoint building failed: %w", err)
	}

	return e.DeleteResource(endpoint)
}
