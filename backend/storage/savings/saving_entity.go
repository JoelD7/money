package savings

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"time"
)

var (
	// This value indicates that a saving hasn't a saving goal associated. It cannot be left empty because saving_goal_id
	// belongs to one of the indices of the table.
	savingGoalIDNone = "none"
)

type savingEntity struct {
	SavingID            string    `json:"saving_id,omitempty"  dynamodbav:"saving_id"`
	SavingGoalID        *string   `json:"saving_goal_id,omitempty"  dynamodbav:"saving_goal_id"`
	Username            string    `json:"username,omitempty"  dynamodbav:"username"`
	PeriodID            *string   `json:"period_id,omitempty"  dynamodbav:"period_id,omitempty"`
	PeriodUser          *string   `json:"period_user,omitempty"  dynamodbav:"period_user"`
	CreatedDate         time.Time `json:"created_date,omitempty"  dynamodbav:"created_date"`
	UpdatedDate         time.Time `json:"updated_date,omitempty"  dynamodbav:"updated_date"`
	Amount              *float64  `json:"amount" dynamodbav:"amount"`
	CreatedDateSavingID string    `json:"created_date_saving_id,omitempty" dynamodbav:"created_date_saving_id"`
}

func toSavingEntity(s *models.Saving) *savingEntity {
	savingEnt := &savingEntity{
		SavingID:     s.SavingID,
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		PeriodID:     s.PeriodID,
		PeriodUser:   s.PeriodUser,
		CreatedDate:  s.CreatedDate,
		UpdatedDate:  s.UpdatedDate,
		Amount:       s.Amount,
		CreatedDateSavingID: dynamo.BuildCreatedDateEntityIDKey(
			s.CreatedDate,
			s.SavingID,
		),
	}

	if s.SavingGoalID == nil || (s.SavingGoalID != nil && *s.SavingGoalID == "") {
		savingEnt.SavingGoalID = &savingGoalIDNone
	}

	return savingEnt
}

func toSavingModel(s savingEntity) *models.Saving {
	savingModel := &models.Saving{
		SavingID:     s.SavingID,
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		PeriodID:     s.PeriodID,
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

func toSavingModels(savings []savingEntity) []*models.Saving {
	modelSavings := make([]*models.Saving, 0, len(savings))

	for _, v := range savings {
		modelSavings = append(modelSavings, toSavingModel(v))
	}

	return modelSavings
}
