package expenses_recurring

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type ExpenseRecurringEntity struct {
	ID           string    `json:"id" dynamodbav:"id"`
	Username     string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID   *string   `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount       float64   `json:"amount" dynamodbav:"amount"`
	Name         string    `json:"name,omitempty" dynamodbav:"name"`
	RecurringDay int       `json:"recurring_day,omitempty" dynamodbav:"recurring_day"`
	Notes        string    `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate  time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdateDate   time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
}

func toExpenseRecurringEntity(e *models.ExpenseRecurring) *ExpenseRecurringEntity {
	return &ExpenseRecurringEntity{
		ID:           e.ID,
		Username:     e.Username,
		CategoryID:   e.CategoryID,
		Amount:       e.Amount,
		Name:         e.Name,
		RecurringDay: e.RecurringDay,
		Notes:        e.Notes,
		CreatedDate:  e.CreatedDate,
		UpdateDate:   e.UpdateDate,
	}
}

func toExpenseRecurringModel(e ExpenseRecurringEntity) *models.ExpenseRecurring {
	return &models.ExpenseRecurring{
		ID:           e.ID,
		Username:     e.Username,
		CategoryID:   e.CategoryID,
		Amount:       e.Amount,
		Name:         e.Name,
		RecurringDay: e.RecurringDay,
		Notes:        e.Notes,
		CreatedDate:  e.CreatedDate,
		UpdateDate:   e.UpdateDate,
	}
}

// Currently the mapping from entity to model is 1:1 so unmarshalling the query result directly to the model might seem
// the obvious option. However, we don't know if the model or entity will change in the future, so it's better to keep
// the mapping logic separate.
func toExpensesRecurringModel(entities []*ExpenseRecurringEntity) []*models.ExpenseRecurring {
	expenses := make([]*models.ExpenseRecurring, 0, len(entities))
	for _, e := range entities {
		expenses = append(expenses, toExpenseRecurringModel(*e))
	}

	return expenses
}

func (e *ExpenseRecurringEntity) Key() string {
	return "expense_entity"
}

func (e *ExpenseRecurringEntity) Value() map[string]interface{} {
	return map[string]interface{}{
		"id":            e.ID,
		"username":      e.Username,
		"category_id":   e.CategoryID,
		"amount":        e.Amount,
		"recurring_day": e.RecurringDay,
		"name":          e.Name,
		"notes":         e.Notes,
		"created_date":  e.CreatedDate,
		"update_date":   e.UpdateDate,
	}
}
