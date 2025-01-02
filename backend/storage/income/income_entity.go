package income

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"time"
)

type incomeEntity struct {
	Username    string    `json:"username,omitempty" dynamodbav:"username"`
	IncomeID    string    `json:"income_id,omitempty" dynamodbav:"income_id"`
	Amount      *float64  `json:"amount" dynamodbav:"amount"`
	Name        *string   `json:"name,omitempty" dynamodbav:"name"`
	Notes       *string   `json:"notes,omitempty" dynamodbav:"notes"`
	CreatedDate time.Time `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty" dynamodbav:"updated_date"`
	Period      *string   `json:"period,omitempty" dynamodbav:"period"`
	PeriodUser  *string   `json:"period_user,omitempty" dynamodbav:"period_user"`
	// AmountKey is a special attribute used to sort income by amount. It's composed of a padded-string of the amount
	AmountKey string `json:"amount_key,omitempty" dynamodbav:"amount_key"`
	// NameIncomeID is a special attribute used to sort income by name. It's composed of the name plus the income id.
	NameIncomeID string `json:"name_income_id,omitempty" dynamodbav:"name_income_id"`
}

func toIncomeModel(i incomeEntity) *models.Income {
	return &models.Income{
		Username:    i.Username,
		IncomeID:    i.IncomeID,
		Amount:      i.Amount,
		Name:        i.Name,
		CreatedDate: i.CreatedDate,
		UpdatedDate: i.UpdatedDate,
		Period:      i.Period,
		Notes:       i.Notes,
		PeriodUser:  i.PeriodUser,
	}
}

func toIncomeModels(is []incomeEntity) []*models.Income {
	incomes := make([]*models.Income, len(is))
	for i, e := range is {
		incomes[i] = toIncomeModel(e)
	}

	return incomes
}

func toIncomeEntity(i *models.Income) *incomeEntity {
	return &incomeEntity{
		Username:     i.Username,
		IncomeID:     i.IncomeID,
		Amount:       i.Amount,
		Name:         i.Name,
		CreatedDate:  i.CreatedDate,
		UpdatedDate:  i.UpdatedDate,
		Period:       i.Period,
		Notes:        i.Notes,
		PeriodUser:   i.PeriodUser,
		AmountKey:    dynamo.BuildAmountKey(*i.Amount, i.IncomeID),
		NameIncomeID: dynamo.BuildNameKey(*i.Name, i.IncomeID),
	}
}
