package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"io"
	"net/http"
)

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
