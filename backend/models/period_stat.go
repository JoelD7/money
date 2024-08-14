package models

type PeriodStat struct {
	PeriodUser string  `json:"period_user"`
	CategoryID string  `json:"category_id"`
	Total      float64 `json:"total"`
}
