package models

import "time"

type InvalidToken struct {
	Email       string    `json:"email,omitempty" dynamodbav:"email"`
	Token       string    `json:"token,omitempty" dynamodbav:"token"`
	Expire      int64     `json:"expire,omitempty" dynamodbav:"expire"`
	Type        string    `json:"type,omitempty" dynamodbav:"type"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date,omitempty"`
}
