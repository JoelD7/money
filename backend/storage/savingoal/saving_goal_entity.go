package savingoal

import (
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
)

type savingGoalEntity struct {
	SavingGoalID     string     `json:"saving_goal_id,omitempty" dynamodbav:"saving_goal_id"`
	Username         string     `json:"username,omitempty" dynamodbav:"username"`
	Name             string     `json:"name,omitempty" dynamodbav:"name"`
	Target           float64    `json:"target,omitempty" dynamodbav:"target"`
	Deadline         time.Time  `json:"deadline,omitempty" dynamodbav:"deadline"`
	CreatedAt        *time.Time `json:"created_at,omitempty" dynamodbav:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" dynamodbav:"updated_at"`
	IsRecurring      bool       `json:"is_recurring,omitempty" dynamodbav:"is_recurring"`
	RecurringAmount  float64    `json:"recurring_amount,omitempty" dynamodbav:"recurring_amount"`
	NameSavingGoalID string     `json:"name-saving_goal_id,omitempty" dynamodbav:"name-saving_goal_id"`
}

func toSavingGoalEntity(s *models.SavingGoal) *savingGoalEntity {
	return &savingGoalEntity{
		SavingGoalID:     s.GetSavingGoalID(),
		Username:         s.GetUsername(),
		Name:             s.GetName(),
		Target:           s.GetTarget(),
		Deadline:         s.GetDeadline(),
		IsRecurring:      s.GetIsRecurring(),
		RecurringAmount:  s.GetRecurringAmount(),
		NameSavingGoalID: dynamo.BuildNameKey(s.GetName(), s.GetSavingGoalID()),
	}
}

func toSavingGoalModel(s *savingGoalEntity) *models.SavingGoal {
	namePtr := &s.Name
	deadlinePtr := &s.Deadline
	targetPtr := &s.Target
	recurringAmountPtr := &s.RecurringAmount

	return &models.SavingGoal{
		SavingGoalID:    s.SavingGoalID,
		Username:        s.Username,
		Name:            namePtr,
		Target:          targetPtr,
		Deadline:        deadlinePtr,
		IsRecurring:     s.IsRecurring,
		RecurringAmount: recurringAmountPtr,
	}
}

func toSavingGoalModels(entities []*savingGoalEntity) []*models.SavingGoal {
	savingGoals := make([]*models.SavingGoal, 0, len(entities))

	for _, savingGoal := range entities {
		savingGoals = append(savingGoals, toSavingGoalModel(savingGoal))
	}

	return savingGoals
}
