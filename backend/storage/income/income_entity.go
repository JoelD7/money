package income

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type incomeEntity struct {
	Username string    `json:"username,omitempty" dynamodbav:"username"`
	IncomeID string    `json:"income_id,omitempty" dynamodbav:"income_id"`
	Amount   float64   `json:"amount" dynamodbav:"amount"`
	Name     string    `json:"name,omitempty" dynamodbav:"name"`
	Date     time.Time `json:"date,omitempty" dynamodbav:"date"`
	Period   string    `json:"period,omitempty" dynamodbav:"period"`
}

func toIncomeModel(i *incomeEntity) *models.Income {
	return &models.Income{
		Username: i.Username,
		IncomeID: i.IncomeID,
		Amount:   i.Amount,
		Name:     i.Name,
		Date:     i.Date,
		Period:   i.Period,
	}
}

func toIncomeModels(is []*incomeEntity) []*models.Income {
	incomes := make([]*models.Income, len(is))
	for i, e := range is {
		incomes[i] = toIncomeModel(e)
	}

	return incomes
}