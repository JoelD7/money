package models

import "time"

type Saving struct {
	SavingID       string    `json:"saving_id,omitempty"  dynamodbav:"saving_id"`
	SavingGoalID   string    `json:"saving_goal_id,omitempty"  dynamodbav:"saving_goal_id"`
	SavingGoalName string    `json:"saving_goal_name,omitempty"  dynamodbav:"saving_goal_name"`
	Username       string    `json:"username,omitempty"  dynamodbav:"username"`
	Period         string    `json:"period,omitempty"  dynamodbav:"period"`
	PeriodUser     string    `json:"period_user,omitempty"  dynamodbav:"period_user"`
	CreatedDate    time.Time `json:"created_date,omitempty"  dynamodbav:"created_date"`
	UpdatedDate    time.Time `json:"updated_date,omitempty"  dynamodbav:"updated_date"`
	Amount         float64   `json:"amount" dynamodbav:"amount"`
}
