package period_stat

import (
	"github.com/JoelD7/money/backend/models"
)

type periodStatEntity struct {
	PeriodUser       string  `json:"period_user" dynamodbav:"period_user"`
	CategoryID       string  `json:"category_id" dynamodbav:"category_id"`
	Total            float64 `json:"total" dynamodbav:"total"`
	CategoryUsername string  `json:"category_username" dynamodbav:"category_username"`
}

func toPeriodStatModel(p periodStatEntity) *models.PeriodStat {
	return &models.PeriodStat{
		PeriodUser: p.PeriodUser,
		CategoryID: p.CategoryID,
		Total:      p.Total,
	}
}
