package models

import "time"

type Income struct {
	Username    string    `json:"username,omitempty"`
	IncomeID    string    `json:"income_id,omitempty"`
	Amount      *float64  `json:"amount"`
	Name        *string   `json:"name,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
	PeriodID    *string   `json:"period_id,omitempty"`
	PeriodUser  *string   `json:"period_user,omitempty"`
	PeriodName  string    `json:"period_name,omitempty"`
}

func (i *Income) GetPeriodID() string {
	if i.PeriodID != nil {
		return *i.PeriodID
	}

	return ""
}

func (i *Income) SetPeriodName(name string) {
	i.PeriodName = name
}

func (i *Income) GetUsername() string {
	return i.Username
}

func (i *Income) GetName() string {
	if i.Name != nil {
		return *i.Name
	}

	return ""
}

func (i *Income) GetAmount() float64 {
	if i.Amount != nil {
		return *i.Amount
	}

	return 0
}
