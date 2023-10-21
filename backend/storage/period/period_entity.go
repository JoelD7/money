package period

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type periodEntity struct {
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	ID          string    `json:"period_id,omitempty" dynamodbav:"period_id"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	StartDate   time.Time `json:"start_date,omitempty" dynamodbav:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty" dynamodbav:"end_date"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty" dynamodbav:"updated_date"`
}

func toPeriodModel(p *periodEntity) *models.Period {
	return &models.Period{
		Username:    p.Username,
		ID:          p.ID,
		Name:        p.Name,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		CreatedDate: p.CreatedDate,
		UpdatedDate: p.UpdatedDate,
	}
}

func toPeriodModels(periods []*periodEntity) []*models.Period {
	periodModels := make([]*models.Period, 0, len(periods))

	for _, period := range periods {
		periodModels = append(periodModels, toPeriodModel(period))
	}

	return periodModels
}

func toPeriodEntity(period *models.Period) *periodEntity {
	return &periodEntity{
		Username:    period.Username,
		ID:          period.ID,
		Name:        period.Name,
		StartDate:   period.StartDate,
		EndDate:     period.EndDate,
		CreatedDate: period.CreatedDate,
		UpdatedDate: period.UpdatedDate,
	}
}
