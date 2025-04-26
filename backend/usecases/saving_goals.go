package usecases

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
	"sync"
)

func NewSavingGoalCreator(savingGoalManager SavingGoalManager) func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	return func(ctx context.Context, username string, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
		savingGoal.Username = username
		return savingGoalManager.CreateSavingGoal(ctx, savingGoal)
	}
}

func NewSavingGoalGetter(savingGoalManager SavingGoalManager, savingManager SavingsManager) func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
	return func(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
		savingGoal, err := savingGoalManager.GetSavingGoal(ctx, username, savingGoalID)
		if err != nil {
			return nil, err
		}

		calculateProgressByGoal(ctx, savingGoal, savingManager)

		return savingGoal, nil
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
				defer func() {
					wg.Done()
				}()
				calculateProgressByGoal(ctx, savingGoal, savingManager)
			}(goal)
		}

		wg.Wait()

		return savingGoals, nextKey, nil
	}
}

func calculateProgressByGoal(ctx context.Context, savingGoal *models.SavingGoal, savingManager SavingsManager) {
	params := &models.QueryParameters{
		PageSize:     10,
		SavingGoalID: savingGoal.SavingGoalID,
	}

	goalSavings := make([]*models.Saving, 0)

	for {
		savings, nextKey, err := savingManager.GetSavingsBySavingGoal(ctx, params)
		if errors.Is(err, models.ErrSavingsNotFound) {
			logger.Info("saving_goal_has_no_savings", models.Any("saving_goal", savingGoal))
			savingGoal.SetProgress(0)
			return
		}

		if errors.Is(err, models.ErrNoMoreItemsToBeRetrieved) {
			logger.Error("no_more_savings_to_be_retrieved", err, models.Any("saving_goal", savingGoal))
			break
		}

		if err != nil {
			logger.Error("calculate_saving_progress_by_goal_failed", err, models.Any("saving_goal", savingGoal),
				models.Any("start_key", params.StartKey))
			savingGoal.SetProgress(0)
			return
		}

		params.StartKey = nextKey

		goalSavings = append(goalSavings, savings...)

		if nextKey == "" {
			break
		}
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
