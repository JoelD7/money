package models

import (
	"time"
)

type ExpenseRecurring struct {
	ID           string    `json:"id"`
	Username     string    `json:"username,omitempty"`
	CategoryID   *string   `json:"category_id,omitempty"`
	Amount       float64   `json:"amount"`
	RecurringDay int       `json:"recurring_day,omitempty"`
	Name         string    `json:"name,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedDate  time.Time `json:"created_date,omitempty"`
	UpdateDate   time.Time `json:"update_date,omitempty"`
}
