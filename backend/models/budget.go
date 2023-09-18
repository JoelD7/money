package models

type Budget struct {
	Username string  `json:"username,omitempty"`
	Month    int32   `json:"month,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
}
