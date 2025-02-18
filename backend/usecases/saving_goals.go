package usecases

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
	"sync"
)

type SavingGoalManager interface {
	CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	UpdateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error)
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error)
	DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error
}

func NewSavingGoalCreator(savingGoalManager SavingGoalManager) func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		return savingGoalManager.CreateSavingGoal(ctx, savingGoal)
	}
}

func NewSavingGoalGetter(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
	return func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
		return savingGoalManager.GetSavingGoal(ctx, username, savingGoalID)
	}
}

func NewSavingGoalsGetter(savingGoalManager SavingGoalManager, savingManager SavingsManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
		savingGoals, nextKey, err := savingGoalManager.GetSavingGoals(ctx, username, params)
		if err != nil {
			return nil, "", err
		}

		var wg sync.WaitGroup

		for _, goal := range savingGoals {
			wg.Add(1)
			go func(savingGoal *models.SavingGoal) {
				defer wg.Done()
				calculateProgressByGoal(ctx, savingGoal, savingManager)
			}(goal)
		}

		wg.Wait()

		return savingGoals, nextKey, nil
	}
}

func calculateProgressByGoal(ctx context.Context, savingGoal *models.SavingGoal, savingManager SavingsManager) {
	startKey := ""
	goalSavings := make([]*models.Saving, 0)

	for {
		savings, nextKey, err := savingManager.GetSavingsBySavingGoal(ctx, startKey, savingGoal.SavingGoalID, 10)
		if errors.Is(err, models.ErrSavingsNotFound) {
			logger.Info("saving_goal_has_no_savings", models.Any("saving_goal", savingGoal))
			savingGoal.SetProgress(0)
			return
		}

		if errors.Is(err, models.ErrNoMoreItemsToBeRetrieved) {
			break
		}

		if err != nil {
			logger.Error("calculate_saving_progress_by_goal_failed", err, models.Any("saving_goal", savingGoal),
				models.Any("start_key", startKey))
			savingGoal.SetProgress(0)
			return
		}

		startKey = nextKey

		goalSavings = append(goalSavings, savings...)
	}

	progress := 0.0
	for _, saving := range goalSavings {
		progress += saving.GetAmount()
	}

	savingGoal.SetProgress(progress)
}

func NewSavingGoalUpdator(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username, savingGoalID string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		savingGoal.SavingGoalID = savingGoalID
		return savingGoalManager.UpdateSavingGoal(ctx, savingGoal)
	}
}

func NewSavingGoalEliminator(savingGoalManager SavingGoalManager) func(ctx context.Context, username, savingGoalID string) error {
	return func(ctx context.Context, username, savingGoalID string) error {
		return savingGoalManager.DeleteSavingGoal(ctx, username, savingGoalID)
	}
}
