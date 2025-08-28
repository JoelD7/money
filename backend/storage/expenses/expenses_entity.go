package expenses

import (
	"strings"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
	er "github.com/JoelD7/money/backend/storage/expenses-recurring"
)

type expenseEntity struct {
	ExpenseID   string    `json:"expense_id" dynamodbav:"expense_id"`
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID  *string   `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	Notes       string    `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	PeriodID    string    `json:"period_id,omitempty" dynamodbav:"period_id"`
	PeriodUser  *string   `json:"period_user,omitempty" dynamodbav:"period_user"`
	UpdateDate  time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
	// AmountKey is a special attribute used to sort expenses by amount. It's composed of a padded-string of the amount
	// plus the expense id.
	AmountKey string `json:"amount_key,omitempty" dynamodbav:"amount_key"`
	// NameExpenseID is a special attribute used to sort expenses by name. It's composed of the name plus the expense id.
	NameExpenseID string `json:"name_expense_id,omitempty" dynamodbav:"name_expense_id"`
}

func toExpenseEntity(e *models.Expense) *expenseEntity {
	entity := &expenseEntity{
		ExpenseID:     e.ExpenseID,
		Username:      e.Username,
		CategoryID:    e.CategoryID,
		Notes:         e.Notes,
		CreatedDate:   e.CreatedDate,
		PeriodID:      e.PeriodID,
		UpdateDate:    e.UpdateDate,
		AmountKey:     dynamo.BuildAmountKey(e.GetAmount(), e.ExpenseID),
		NameExpenseID: dynamo.BuildNameKey(e.GetName(), e.ExpenseID),
	}

	if e.Amount != nil {
		entity.Amount = *e.Amount
	}

	if e.Name != nil {
		entity.Name = *e.Name
	}

	return entity
}

func toExpenseModel(e expenseEntity) *models.Expense {
	return &models.Expense{
		ExpenseID:   e.ExpenseID,
		Username:    e.Username,
		CategoryID:  e.CategoryID,
		Amount:      &e.Amount,
		Name:        &e.Name,
		Notes:       e.Notes,
		CreatedDate: e.CreatedDate,
		PeriodID:    e.PeriodID,
		UpdateDate:  e.UpdateDate,
	}
}

func toExpenseModels(es []expenseEntity) []*models.Expense {
	expenses := make([]*models.Expense, len(es))
	for i, e := range es {
		expenses[i] = toExpenseModel(e)
	}

	return expenses
}

func toExpenseRecurringEntity(e *models.Expense) *er.ExpenseRecurringEntity {
	entity := &er.ExpenseRecurringEntity{
		ID:           strings.ToLower(*e.Name),
		Username:     e.Username,
		CategoryID:   e.CategoryID,
		Notes:        e.Notes,
		RecurringDay: *e.RecurringDay,
		CreatedDate:  e.CreatedDate,
		UpdateDate:   e.UpdateDate,
	}

	if e.Amount != nil {
		entity.Amount = *e.Amount
	}

	if e.Name != nil {
		entity.Name = *e.Name
	}

	return entity
}

func (e *expenseEntity) Key() string {
	return "expense_entity"
}

func (e *expenseEntity) Value() map[string]interface{} {
	return map[string]interface{}{
		"expense_id":   e.ExpenseID,
		"username":     e.Username,
		"category_id":  e.CategoryID,
		"amount":       e.Amount,
		"name":         e.Name,
		"notes":        e.Notes,
		"created_date": e.CreatedDate,
		"period_id":    e.PeriodID,
		"update_date":  e.UpdateDate,
	}
}
