package models

import "time"

type SavingGoal struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	Name         string    `json:"name,omitempty"`
	Goal         float64   `json:"goal,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty"`
}
