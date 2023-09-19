package models

import "time"

type Saving struct {
	SavingID     string    `json:"saving_id,omitempty"  dynamodbav:"saving_id"`
	SavingGoalID string    `json:"saving_goal_id,omitempty"  dynamodbav:"saving_goal_id"`
	Username     string    `json:"username,omitempty"  dynamodbav:"username"`
	CreatedDate  time.Time `json:"created_date,omitempty"  dynamodbav:"created_date"`
	UpdatedDate  time.Time `json:"updated_date,omitempty"  dynamodbav:"updated_date"`
	Amount       float64   `json:"amount" dynamodbav:"amount"`
}
