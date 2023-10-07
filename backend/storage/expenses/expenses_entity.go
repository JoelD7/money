package expenses

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type expenseEntity struct {
	ExpenseID  string    `json:"expense_id" dynamodbav:"expense_id"`
	Username   string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID string    `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount     float64   `json:"amount" dynamodbav:"amount"`
	Currency   string    `json:"currency,omitempty" dynamodbav:"currency"`
	Name       string    `json:"name,omitempty" dynamodbav:"name"`
	Notes      string    `json:"notes,omitempty" dynamodbav:"notes"`
	Date       time.Time `json:"date,omitempty" dynamodbav:"date"`
	Period     string    `json:"period,omitempty" dynamodbav:"period"`
	UpdateDate time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
}

func toExpenseModel(e *expenseEntity) *models.Expense {
	return &models.Expense{
		ExpenseID:  e.ExpenseID,
		Username:   e.Username,
		CategoryID: e.CategoryID,
		Amount:     e.Amount,
		Currency:   e.Currency,
		Name:       e.Name,
		Notes:      e.Notes,
		Date:       e.Date,
		Period:     e.Period,
		UpdateDate: e.UpdateDate,
	}
}

func toExpenseModels(es []*expenseEntity) []*models.Expense {
	expenses := make([]*models.Expense, len(es))
	for i, e := range es {
		expenses[i] = toExpenseModel(e)
	}

	return expenses
}
