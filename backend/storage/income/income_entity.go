package income

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type incomeEntity struct {
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	IncomeID    string    `json:"income_id,omitempty" dynamodbav:"income_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Name        string    `json:"name,omitempty" dynamodbav:"name"`
	Notes       string    `json:"notes,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty"`
	Period      string    `json:"period,omitempty" dynamodbav:"period"`
}

func toIncomeModel(i *incomeEntity) *models.Income {
	return &models.Income{
		Username:    i.Username,
		IncomeID:    i.IncomeID,
		Amount:      i.Amount,
		Name:        i.Name,
		CreatedDate: i.CreatedDate,
		UpdatedDate: i.UpdatedDate,
		Period:      i.Period,
		Notes:       i.Notes,
	}
}

func toIncomeModels(is []*incomeEntity) []*models.Income {
	incomes := make([]*models.Income, len(is))
	for i, e := range is {
		incomes[i] = toIncomeModel(e)
	}

	return incomes
}
