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

func (e *Expense) GetPeriodID() string {
	return e.PeriodID
}

func (e *Expense) SetPeriodName(name string) {
	e.PeriodName = name
}

func (e *Expense) GetUsername() string {
	return e.Username
}

func (e *Expense) GetName() string {
	if e.Name != nil {
		return *e.Name
	}

	return ""
}

func (e *Expense) GetAmount() float64 {
	if e.Amount != nil {
		return *e.Amount
	}

	return 0
}
