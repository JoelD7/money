package entities

import "time"

type SavingGoal struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty" dynamodbav:"saving_goal_id,omitempty"`
	Username     string    `json:"username,omitempty" dynamodbav:"username,omitempty"`
	Name         string    `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Goal         float64   `json:"goal,omitempty" dynamodbav:"goal,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty" dynamodbav:"deadline,omitempty"`
}
