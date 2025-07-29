package models

import "time"

type Expense struct {
	ExpenseID    string    `json:"expense_id"`
	Username     string    `json:"username,omitempty"`
	CategoryID   *string   `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	Amount       *float64  `json:"amount"`
	RecurringDay *int      `json:"recurring_day,omitempty"`
	IsRecurring  bool      `json:"is_recurring"`
	Name         *string   `json:"name,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedDate  time.Time `json:"created_date,omitempty"`
	PeriodID     string    `json:"period_id,omitempty"`
	PeriodName   string    `json:"period_name,omitempty"`
	PeriodUser   *string   `json:"period_user,omitempty"`
	UpdateDate   time.Time `json:"update_date,omitempty"`
}

type CategoryExpenseSummary struct {
	CategoryID string  `json:"category_id"`
	Total      float64 `json:"total"`
	Period     string  `json:"period,omitempty"`
}
