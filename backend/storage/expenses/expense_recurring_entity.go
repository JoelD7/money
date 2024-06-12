package expenses

import (
	"github.com/JoelD7/money/backend/models"
	"strings"
	"time"
)

type expenseRecurringEntity struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID  *string   `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	Notes       string    `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdateDate  time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
}

func toExpenseRecurringEntity(e *models.Expense) *expenseRecurringEntity {
	entity := &expenseRecurringEntity{
		ID:          strings.ToLower(*e.Name),
		Username:    e.Username,
		CategoryID:  e.CategoryID,
		Notes:       e.Notes,
		CreatedDate: e.CreatedDate,
		UpdateDate:  e.UpdateDate,
	}

	if e.Amount != nil {
		entity.Amount = *e.Amount
	}

	if e.Name != nil {
		entity.Name = *e.Name
	}

	return entity
}
