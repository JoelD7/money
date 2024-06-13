package expenses_recurring

import (
	"time"
)

type ExpenseRecurringEntity struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	CategoryID  *string   `json:"category_id,omitempty" dynamodbav:"category_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	Notes       string    `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdateDate  time.Time `json:"update_date,omitempty" dynamodbav:"update_date"`
}
