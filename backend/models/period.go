package models

import "time"

type Period struct {
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	ID          string    `json:"period_id,omitempty" dynamodbav:"period_id"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty" dynamodbav:"updated_date"`
}
