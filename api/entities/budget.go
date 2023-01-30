package entities

type Budget struct {
	PersonID string  `json:"person_id,omitempty"`
	Month    int32   `json:"month,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
}
