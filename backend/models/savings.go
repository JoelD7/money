package models

import "time"

type Saving struct {
	SavingID     string    `json:"saving_id,omitempty"`
	SavingGoalID string    `json:"saving_goal_id,omitempty"`
	Email        string    `json:"email,omitempty"`
	CreationDate time.Time `json:"creation_date,omitempty"`
	Amount       float64   `json:"amount"`
}
