package models

type InvalidToken struct {
	Email  string `json:"email,omitempty" dynamodbav:"email"`
	Token  string `json:"token,omitempty" dynamodbav:"token"`
	Expire int64  `json:"expire,omitempty" dynamodbav:"expire"`
}
