package models

import "time"

type Saving struct {
	SavingID     string    `json:"saving_id,omitempty"  dynamodbav:"saving_id"`
	SavingGoalID string    `json:"saving_goal_id,omitempty"  dynamodbav:"saving_goal_id"`
	Email        string    `json:"email,omitempty"  dynamodbav:"email"`
	CreationDate time.Time `json:"creation_date,omitempty"  dynamodbav:"creation_date"`
	Amount       float64   `json:"amount" dynamodbav:"amount"`
}
