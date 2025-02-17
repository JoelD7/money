package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"io"
	"net/http"
)

type savingGoalsResponse struct {
	SavingGoals []*models.SavingGoal `json:"saving_goals"`
	NextKey     string               `json:"next_key"`
}

func (e *E2ERequester) CreateSavingGoal(savingGoal *models.SavingGoal) (*models.SavingGoal, int, error) {
	requestBody, err := json.Marshal(savingGoal)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, e.baseUrl+savingGoalsEndpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request failed: %w", err)
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

	var createdSavingGoal models.SavingGoal
	err = json.NewDecoder(res.Body).Decode(&createdSavingGoal)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("saving goal response decoding failed: %w", err)
	}

	return &createdSavingGoal, res.StatusCode, nil
}

func (e *E2ERequester) GetSavingGoal(savingGoalID string) (*models.SavingGoal, int, error) {
	request, err := http.NewRequest(http.MethodGet, e.baseUrl+savingGoalsEndpoint+"/"+savingGoalID, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	var savingGoal models.SavingGoal
	err = json.NewDecoder(res.Body).Decode(&savingGoal)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("saving goal response decoding failed: %w", err)
	}

	return &savingGoal, res.StatusCode, nil
}

func (e *E2ERequester) DeleteSavingGoal(savingGoalID string) (int, error) {
	request, err := http.NewRequest(http.MethodDelete, e.baseUrl+savingGoalsEndpoint+"/"+savingGoalID, nil)
	if err != nil {
		return 0, fmt.Errorf("saving goal request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("saving goal request failed: %w", err)
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

func (e *E2ERequester) GetSavingGoals(sortBy, sortOrder, startKey string, pageSize int) ([]*models.SavingGoal, int, string, error) {
	request, err := http.NewRequest(http.MethodGet, e.baseUrl+savingGoalsEndpoint, nil)
	if err != nil {
		return nil, 0, "", fmt.Errorf("saving goals request building failed: %w", err)
	}

	q := request.URL.Query()
	q.Add("sort_by", sortBy)
	q.Add("sort_order", sortOrder)
	q.Add("start_key", startKey)
	q.Add("page_size", fmt.Sprint(pageSize))
	request.URL.RawQuery = q.Encode()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("saving goals request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, "", handleErrorResponse(res.StatusCode, res.Body)
	}

	var response savingGoalsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, res.StatusCode, "", fmt.Errorf("saving goals response decoding failed: %w", err)
	}

	return response.SavingGoals, res.StatusCode, response.NextKey, nil
}

func (e *E2ERequester) UpdateSavingGoal(savingGoalID string, savingGoal *models.SavingGoal) (*models.SavingGoal, int, error) {
	requestBody, err := json.Marshal(savingGoal)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request body marshalling failed: %w", err)
	}

	request, err := http.NewRequest(http.MethodPut, e.baseUrl+savingGoalsEndpoint+"/"+savingGoalID, bytes.NewReader(requestBody))
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request building failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Auth", "Bearer "+e.accessToken)

	res, err := e.client.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("saving goal request failed: %w", err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	var updatedSavingGoal models.SavingGoal
	err = json.NewDecoder(res.Body).Decode(&updatedSavingGoal)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("saving goal response decoding failed: %w", err)
	}

	return &updatedSavingGoal, res.StatusCode, nil
}

func handleErrorResponse(statusCode int, body io.ReadCloser) error {
	var errRes ErrorResponse
	err := json.NewDecoder(body).Decode(&errRes)
	if err != nil {
		return fmt.Errorf("error response decoding failed: %w", err)
	}

	if errRes.Message == "" {
		return fmt.Errorf("saving goal request failed with status: %v", statusCode)
	}

	return fmt.Errorf(errRes.Message)
}
