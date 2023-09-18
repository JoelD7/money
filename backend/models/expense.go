package models

import "time"

type Expense struct {
	ExpenseID     string    `json:"expense_id"`
	Username      string    `json:"username,omitempty"`
	CategoryID    string    `json:"category_id,omitempty"`
	SubcategoryID string    `json:"subcategory_id,omitempty"`
	SavingGoalID  string    `json:"saving_goal_id,omitempty"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency,omitempty"`
	Name          string    `json:"name,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	Date          time.Time `json:"date,omitempty"`
	Period        string    `json:"period,omitempty"`
	UpdateDate    time.Time `json:"update_date,omitempty"`
}
