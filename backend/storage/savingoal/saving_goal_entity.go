package savingoal

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type savingGoalEntity struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty" dynamodbav:"saving_goal_id"`
	Username     string    `json:"username,omitempty" dynamodbav:"username"`
	Name         string    `json:"name,omitempty" dynamodbav:"name"`
	Target       float64   `json:"target,omitempty" dynamodbav:"target"`
	Deadline     time.Time `json:"deadline,omitempty" dynamodbav:"deadline"`
}

func toSavingGoalEntity(s *models.SavingGoal) *savingGoalEntity {
	return &savingGoalEntity{
		SavingGoalID: s.GetSavingGoalID(),
		Username:     s.GetUsername(),
		Name:         s.GetName(),
		Target:       s.GetTarget(),
		Deadline:     s.GetDeadline(),
	}
}

func toSavingGoalModel(s *savingGoalEntity) *models.SavingGoal {
	namePtr := &s.Name
	deadlinePtr := &s.Deadline
	targetPtr := &s.Target

	return &models.SavingGoal{
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		Name:         namePtr,
		Target:       targetPtr,
		Deadline:     deadlinePtr,
	}
}

func toSavingGoalModels(entities []*savingGoalEntity) []*models.SavingGoal {
	savingGoals := make([]*models.SavingGoal, 0, len(entities))

	for _, savingGoal := range entities {
		savingGoals = append(savingGoals, toSavingGoalModel(savingGoal))
	}

	return savingGoals
}
