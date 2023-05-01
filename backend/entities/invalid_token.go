package entities

type InvalidToken struct {
	Email  string `json:"email,omitempty" dynamodbav:"email"`
	Token  string `json:"token_id,omitempty" dynamodbav:"token_id"`
	Expire int64  `json:"expire,omitempty" dynamodbav:"expire"`
}
