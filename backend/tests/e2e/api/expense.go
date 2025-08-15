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

type ExpensesResponse struct {
	Expenses []*models.Expense `json:"expenses"`
	NextKey  string            `json:"next_key"`
}

func (e *E2ERequester) CreateExpense(expense *models.Expense, t TestCleaner) (*models.Expense, int, error) {
	var createdExpense models.Expense
	t.Cleanup(func() {
		if createdExpense.ExpenseID == "" {
			return
		}

		statusCode, err := e.DeleteExpense(createdExpense.ExpenseID)
		if statusCode != http.StatusNoContent || err != nil {
			t.Logf("Failed to delete expense %s: %v", createdExpense.ExpenseID, err)
		}
	})

	requestBody, err := json.Marshal(expense)
	if err != nil {
		return nil, 0, fmt.Errorf("request body marshalling failed: %w", err)
	}

	endpoint, err := url.JoinPath(e.baseUrl, expensesEndpoint)
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

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusCreated {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	err = json.NewDecoder(res.Body).Decode(&createdExpense)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("response decoding failed: %w", err)
	}

	return &createdExpense, res.StatusCode, nil
}

func (e *E2ERequester) GetExpenses(params *models.QueryParameters) ([]*models.Expense, int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, expensesEndpoint)
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

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("closing response body failed: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, handleErrorResponse(res.StatusCode, res.Body)
	}

	var expensesRes ExpensesResponse
	err = json.NewDecoder(res.Body).Decode(&expensesRes)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("response decoding failed: %w", err)
	}

	return expensesRes.Expenses, res.StatusCode, nil
}

func addQueryParams(request *http.Request, params *models.QueryParameters, active bool) {
	if params == nil {
		return
	}

	q := request.URL.Query()
	if params.Period != "" {
		q.Add("period", fmt.Sprintf("%s", params.Period))
	}

	if len(params.Categories) > 0 {
		for _, category := range params.Categories {
			q.Add("category", fmt.Sprintf("%s", category))
		}
	}

	if params.SavingGoalID != "" {
		q.Add("saving_goal_id", fmt.Sprintf("%s", params.SavingGoalID))
	}

	if params.PageSize > 0 {
		q.Add("page_size", fmt.Sprintf("%d", params.PageSize))
	}

	if params.SortBy != "" {
		q.Add("sort_by", fmt.Sprintf("%s", params.SortBy))
	}

	if params.SortType != "" {
		q.Add("sort_type", fmt.Sprintf("%s", params.SortType))
	}

	if params.StartKey != "" {
		q.Add("start_key", fmt.Sprintf("%s", params.StartKey))
	}

	if active {
		q.Add("active", "true")
	}

	request.URL.RawQuery = q.Encode()
}

func (e *E2ERequester) DeleteExpense(expenseID string) (int, error) {
	endpoint, err := url.JoinPath(e.baseUrl, expensesEndpoint, expenseID)
	if err != nil {
		return 0, fmt.Errorf("expense deletion endpoint building failed: %w", err)
	}

	return e.DeleteResource(endpoint)
}
