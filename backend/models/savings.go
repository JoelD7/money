package models

import "time"

type Saving struct {
	SavingID       string    `json:"saving_id,omitempty"`
	SavingGoalID   *string   `json:"saving_goal_id,omitempty"`
	SavingGoalName string    `json:"saving_goal_name,omitempty"`
	Username       string    `json:"username,omitempty"`
	PeriodID       *string   `json:"period_id,omitempty"`
	PeriodUser     *string   `json:"-"`
	PeriodName     string    `json:"period_name,omitempty"`
	CreatedDate    time.Time `json:"created_date,omitempty"`
	UpdatedDate    time.Time `json:"updated_date,omitempty"`
	Amount         *float64  `json:"amount"`
}

func (s *Saving) GetPeriodID() string {
	if s.PeriodID != nil {
		return *s.PeriodID
	}

	return ""
}

func (s *Saving) SetPeriodName(name string) {
	s.PeriodName = name
}

func (s *Saving) GetAmount() float64 {
	if s == nil || s.Amount == nil {
		return 0
	}
	return *s.Amount
}

func (s *Saving) GetUsername() string {
	return s.Username
}
