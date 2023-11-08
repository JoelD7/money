package savings

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

var (
	// This value indicates that a saving hasn't a saving goal associated. It cannot be left empty because saving_goal_id
	// belongs to one of the indices of the table.
	savingGoalIDNone = "none"
)

type savingEntity struct {
	SavingID     string    `json:"saving_id,omitempty"  dynamodbav:"saving_id"`
	SavingGoalID *string   `json:"saving_goal_id,omitempty"  dynamodbav:"saving_goal_id"`
	Username     string    `json:"username,omitempty"  dynamodbav:"username"`
	Period       *string   `json:"period,omitempty"  dynamodbav:"period"`
	PeriodUser   *string   `json:"period_user,omitempty"  dynamodbav:"period_user"`
	CreatedDate  time.Time `json:"created_date,omitempty"  dynamodbav:"created_date"`
	UpdatedDate  time.Time `json:"updated_date,omitempty"  dynamodbav:"updated_date"`
	Amount       *float64  `json:"amount" dynamodbav:"amount"`
}

func toSavingEntity(s *models.Saving) *savingEntity {
	savingEnt := &savingEntity{
		SavingID:     s.SavingID,
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		Period:       s.Period,
		PeriodUser:   s.PeriodUser,
		CreatedDate:  s.CreatedDate,
		UpdatedDate:  s.UpdatedDate,
		Amount:       s.Amount,
	}

	if s.SavingGoalID == nil || (s.SavingGoalID != nil && *s.SavingGoalID == "") {
		savingEnt.SavingGoalID = &savingGoalIDNone
	}

	return savingEnt
}

func toSavingModel(s *savingEntity) *models.Saving {
	savingModel := &models.Saving{
		SavingID:     s.SavingID,
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		Period:       s.Period,
		PeriodUser:   s.PeriodUser,
		CreatedDate:  s.CreatedDate,
		UpdatedDate:  s.UpdatedDate,
		Amount:       s.Amount,
	}

	if savingModel.SavingGoalID != nil && *savingModel.SavingGoalID == savingGoalIDNone {
		savingModel.SavingGoalID = nil
	}

	return savingModel
}

func toSavingModels(savings []*savingEntity) []*models.Saving {
	modelSavings := make([]*models.Saving, 0, len(savings))

	for _, v := range savings {
		modelSavings = append(modelSavings, toSavingModel(v))
	}

	return modelSavings
}
