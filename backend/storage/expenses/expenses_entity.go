package expenses

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type expenseEntity struct {
	ExpenseID   string    `json:"expense_id" dynamodbav:"expense_id"`
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID  *string   `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	Notes       string    `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	Period      string    `json:"period,omitempty" dynamodbav:"period"`
	PeriodUser  string    `json:"period_user,omitempty" dynamodbav:"period_user"`
	UpdateDate  time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
}

func toExpenseEntity(e *models.Expense) *expenseEntity {
	return &expenseEntity{
		ExpenseID:   e.ExpenseID,
		Username:    e.Username,
		CategoryID:  e.CategoryID,
		Amount:      *e.Amount,
		Name:        *e.Name,
		Notes:       e.Notes,
		CreatedDate: e.CreatedDate,
		Period:      e.Period,
		UpdateDate:  e.UpdateDate,
	}
}

func toExpenseModel(e *expenseEntity) *models.Expense {
	return &models.Expense{
		ExpenseID:   e.ExpenseID,
		Username:    e.Username,
		CategoryID:  e.CategoryID,
		Amount:      &e.Amount,
		Name:        &e.Name,
		Notes:       e.Notes,
		CreatedDate: e.CreatedDate,
		Period:      e.Period,
		UpdateDate:  e.UpdateDate,
	}
}

func toExpenseModels(es []expenseEntity) []*models.Expense {
	expenses := make([]*models.Expense, len(es))
	for i, e := range es {
		expenses[i] = toExpenseModel(&e)
	}

	return expenses
}
