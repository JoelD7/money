package models

import "time"

type Saving struct {
	SavingID       string    `json:"saving_id,omitempty"`
	SavingGoalID   *string   `json:"saving_goal_id,omitempty"`
	SavingGoalName string    `json:"saving_goal_name,omitempty"`
	Username       string    `json:"username,omitempty"`
	Period         *string   `json:"period,omitempty"`
	PeriodUser     *string   `json:"-"`
	CreatedDate    time.Time `json:"created_date,omitempty"`
	UpdatedDate    time.Time `json:"updated_date,omitempty"`
	Amount         *float64  `json:"amount"`
}
