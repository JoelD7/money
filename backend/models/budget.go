package models

type Budget struct {
	UserID string  `json:"user_id,omitempty"`
	Month  int32   `json:"month,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}
