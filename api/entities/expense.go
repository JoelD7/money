package entities

import "time"

type Expense struct {
	PersonID      string    `json:"person_id,omitempty"`
	CategoryID    string    `json:"category_id,omitempty"`
	SubcategoryID string    `json:"subcategory_id,omitempty"`
	SavingGoalID  string    `json:"saving_goal_id,omitempty"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency,omitempty"`
	Name          string    `json:"name,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	CreationDate  time.Time `json:"creation_date,omitempty"`
	UpdateDate    time.Time `json:"update_date,omitempty"`
}
