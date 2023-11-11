package models

import "time"

type Income struct {
	Username    string    `json:"username,omitempty"`
	IncomeID    string    `json:"income_id,omitempty"`
	Amount      float64   `json:"amount"`
	Name        string    `json:"name,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
	Period      string    `json:"period,omitempty"`
}
