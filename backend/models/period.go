package models

import (
	"time"
)

type PeriodType string

const PeriodTypeCurrent PeriodType = "current"

type Period struct {
	Username    string    `json:"username,omitempty"`
	ID          string    `json:"period,omitempty"`
	Name        *string   `json:"name,omitempty"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
}

func (period *Period) LogName() string {
	return "period"
}

func (period *Period) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_username":     period.Username,
		"s_period_id":    period.ID,
		"s_name":         period.Name,
		"s_start_date":   period.StartDate,
		"s_end_date":     period.EndDate,
		"s_created_date": period.CreatedDate,
		"s_updated_date": period.UpdatedDate,
	}
}

type PeriodStat struct {
	PeriodID               string                    `json:"period_id"`
	TotalIncome            float64                   `json:"total_income"`
	CategoryExpenseSummary []*CategoryExpenseSummary `json:"category_expense_summary"`
}
