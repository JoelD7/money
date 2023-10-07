package savingoal

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type savingGoalEntity struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	Name         string    `json:"name,omitempty"`
	Goal         float64   `json:"goal,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty"`
}

func toSavingGoalModel(s *savingGoalEntity) *models.SavingGoal {
	return &models.SavingGoal{
		SavingGoalID: s.SavingGoalID,
		Username:     s.Username,
		Name:         s.Name,
		Goal:         s.Goal,
		Deadline:     s.Deadline,
	}
}

func toSavingGoalModels(entities []*savingGoalEntity) []*models.SavingGoal {
	savingGoals := make([]*models.SavingGoal, 0, len(entities))

	for _, savingGoal := range entities {
		savingGoals = append(savingGoals, toSavingGoalModel(savingGoal))
	}

	return savingGoals
}
