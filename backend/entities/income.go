package entities

import "time"

type Income struct {
	Username string    `json:"username,omitempty"`
	IncomeID string    `json:"income_id,omitempty"`
	Amount   float64   `json:"amount"`
	Name     string    `json:"name,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Period   string    `json:"period,omitempty"`
}
