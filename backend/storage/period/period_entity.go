package period

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type periodEntity struct {
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	ID          string    `json:"period,omitempty" dynamodbav:"period"`
	Name        *string   `json:"name,omitempty" dynamodbav:"name"`
	StartDate   time.Time `json:"start_date,omitempty" dynamodbav:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty" dynamodbav:"end_date"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty" dynamodbav:"updated_date"`
}

type uniquePeriodNameEntity struct {
	Name     string `json:"name,omitempty" dynamodbav:"name"`
	Username string `json:"username,omitempty" dynamodbav:"username"`
}

func toPeriodModel(p periodEntity) *models.Period {
	start := models.ToPeriodTime(p.StartDate)
	end := models.ToPeriodTime(p.EndDate)

	return &models.Period{
		Username:    p.Username,
		ID:          p.ID,
		Name:        p.Name,
		StartDate:   start,
		EndDate:     end,
		CreatedDate: p.CreatedDate,
		UpdatedDate: p.UpdatedDate,
	}
}

func toPeriodModels(periods []periodEntity) []*models.Period {
	periodModels := make([]*models.Period, 0, len(periods))

	for _, period := range periods {
		periodModels = append(periodModels, toPeriodModel(period))
	}

	return periodModels
}

func toPeriodEntity(period models.Period) periodEntity {
	start := models.ToTime(period.StartDate)
	end := models.ToTime(period.EndDate)

	return periodEntity{
		Username:    period.Username,
		ID:          period.ID,
		Name:        period.Name,
		StartDate:   start,
		EndDate:     end,
		CreatedDate: period.CreatedDate,
		UpdatedDate: period.UpdatedDate,
	}
}
