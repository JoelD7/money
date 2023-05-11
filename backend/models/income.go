package models

import "time"

type Income struct {
	PersonID string    `json:"person_id,omitempty"`
	IncomeID string    `json:"income_id,omitempty"`
	Amount   float64   `json:"amount"`
	Name     string    `json:"name,omitempty"`
	Date     time.Time `json:"date,omitempty"`
}
