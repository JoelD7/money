package models

import "time"

type Expense struct {
	ExpenseID    string    `json:"expense_id"`
	Username     string    `json:"username,omitempty"`
	CategoryID   *string   `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	Amount       *float64  `json:"amount"`
	Name         *string   `json:"name,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedDate  time.Time `json:"created_date,omitempty"`
	Period       *string   `json:"period,omitempty"`
	PeriodUser   *string   `json:"period_user,omitempty"`
	UpdateDate   time.Time `json:"update_date,omitempty"`
}
