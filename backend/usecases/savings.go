package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"math"
	"time"
)

type SavingsManager interface {
	GetSaving(ctx context.Context, username, savingID string) (*models.Saving, error)
	GetSavings(ctx context.Context, username, startKey string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsByPeriod(ctx context.Context, startKey, username, period string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsBySavingGoal(ctx context.Context, startKey, savingGoalID string, pageSize int) ([]*models.Saving, string, error)
	GetSavingsBySavingGoalAndPeriod(ctx context.Context, startKey, savingGoalID, period string, pageSize int) ([]*models.Saving, string, error)
	CreateSaving(ctx context.Context, saving *models.Saving) (*models.Saving, error)
	UpdateSaving(ctx context.Context, saving *models.Saving) error
	DeleteSaving(ctx context.Context, savingID, username string) error
}

type SavingGoalManager interface {
	GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error)
	GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error)
}

func NewSavingGetter(sm SavingsManager, sgm SavingGoalManager, l Logger) func(ctx context.Context, username, savingID string) (*models.Saving, error) {
	return func(ctx context.Context, username, savingID string) (*models.Saving, error) {
		saving, err := sm.GetSaving(ctx, username, savingID)
		if err != nil {
			return nil, err
		}

		err = setSavingGoalName(ctx, sgm, saving)
		if err != nil {
			return saving, fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return saving, nil
	}
}

func NewSavingsGetter(sm SavingsManager, sgm SavingGoalManager, l Logger) func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Saving, string, error) {
	return func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Saving, string, error) {
		if err := validatePageSize(pageSize); err != nil {
			l.Error("invalid_page_size_detected", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"s_username":  username,
					"i_page_size": pageSize,
				}),
			})

			return nil, "", err
		}

		savings, nextKey, err := sm.GetSavings(ctx, username, startKey, pageSize)
		if err != nil {
			return nil, "", fmt.Errorf("savings fetch failed: %w", err)
		}

		err = setSavingGoalNames(ctx, sgm, username, savings)
		if err != nil {
			return savings, "", fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return savings, nextKey, nil
	}
}

func NewSavingByPeriodGetter(sm SavingsManager, sgm SavingGoalManager, l Logger) func(ctx context.Context, username, startKey, period string, pageSize int) ([]*models.Saving, string, error) {
	return func(ctx context.Context, username, startKey, period string, pageSize int) ([]*models.Saving, string, error) {
		if err := validatePageSize(pageSize); err != nil {
			l.Error("invalid_page_size_detected", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"s_username":  username,
					"i_page_size": pageSize,
				}),
			})

			return nil, "", err
		}

		savings, nextKey, err := sm.GetSavingsByPeriod(ctx, startKey, username, period, pageSize)
		if err != nil {
			return nil, "", fmt.Errorf("savings fetch failed: %w", err)
		}

		err = setSavingGoalNames(ctx, sgm, username, savings)
		if err != nil {
			return savings, "", fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return savings, nextKey, nil
	}
}

func NewSavingBySavingGoalGetter(sm SavingsManager, sgm SavingGoalManager, l Logger) func(ctx context.Context, startKey, savingGoalID string, pageSize int) ([]*models.Saving, string, error) {
	return func(ctx context.Context, startKey, savingGoalID string, pageSize int) ([]*models.Saving, string, error) {
		if err := validatePageSize(pageSize); err != nil {
			l.Error("invalid_page_size_detected", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"i_page_size": pageSize,
				}),
			})

			return nil, "", err
		}

		savings, nextKey, err := sm.GetSavingsBySavingGoal(ctx, startKey, savingGoalID, pageSize)
		if err != nil {
			return nil, "", fmt.Errorf("savings fetch failed: %w", err)
		}

		err = setSavingGoalNamesForSavingGoal(ctx, sgm, savings[0].Username, savingGoalID, savings)
		if err != nil {
			return savings, "", fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return savings, nextKey, nil
	}
}

func NewSavingBySavingGoalAndPeriodGetter(sm SavingsManager, sgm SavingGoalManager, l Logger) func(ctx context.Context, startKey, savingGoalID, period string, pageSize int) ([]*models.Saving, string, error) {
	return func(ctx context.Context, startKey, savingGoalID, period string, pageSize int) ([]*models.Saving, string, error) {
		if err := validatePageSize(pageSize); err != nil {
			l.Error("invalid_page_size_detected", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"i_page_size": pageSize,
				}),
			})

			return nil, "", err
		}

		savings, nextKey, err := sm.GetSavingsBySavingGoalAndPeriod(ctx, startKey, savingGoalID, period, pageSize)
		if err != nil {
			return nil, "", fmt.Errorf("savings fetch failed: %w", err)
		}

		err = setSavingGoalNames(ctx, sgm, savings[0].Username, savings)
		if err != nil {
			return savings, "", fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return savings, nextKey, nil
	}
}

func NewSavingCreator(sm SavingsManager, u UserManager, p PeriodManager) func(ctx context.Context, username string, saving *models.Saving) (*models.Saving, error) {
	return func(ctx context.Context, username string, saving *models.Saving) (*models.Saving, error) {
		user, err := u.GetUser(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("user fetch failed: %w", err)
		}

		err = validateSavingPeriod(ctx, saving, username, p)
		if err != nil {
			return nil, err
		}

		if saving.Period == nil || *saving.Period == "" {
			saving.Period = &user.CurrentPeriod
		}

		saving.SavingID = generateDynamoID("SV")
		saving.Username = username
		saving.CreatedDate = time.Now()

		newSaving, err := sm.CreateSaving(ctx, saving)
		if err != nil {
			return nil, fmt.Errorf("saving creation failed: %w", err)
		}

		return newSaving, nil
	}
}

func NewSavingUpdater(sm SavingsManager, pm PeriodManager, sgm SavingGoalManager) func(ctx context.Context, username string, saving *models.Saving) (*models.Saving, error) {
	return func(ctx context.Context, username string, saving *models.Saving) (*models.Saving, error) {
		err := validateSavingPeriod(ctx, saving, username, pm)
		if err != nil {
			return nil, err
		}

		saving.UpdatedDate = time.Now()

		err = sm.UpdateSaving(ctx, saving)
		if err != nil {
			return nil, err
		}

		updatedSaving, err := sm.GetSaving(ctx, username, saving.SavingID)
		if err != nil {
			return nil, fmt.Errorf("getting updated saving failed: %w", err)
		}

		err = setSavingGoalName(ctx, sgm, updatedSaving)
		if err != nil {
			return updatedSaving, fmt.Errorf("%w: %v", models.ErrSavingGoalNameSettingFailed, err)
		}

		return updatedSaving, nil
	}
}

func validatePageSize(pageSize int) error {
	if pageSize < 0 || pageSize > math.MaxInt32 {
		return models.ErrInvalidPageSize
	}

	return nil
}

func NewSavingDeleter(sm SavingsManager) func(ctx context.Context, savingID, username string) error {
	return func(ctx context.Context, savingID, username string) error {
		if savingID == "" {
			return models.ErrMissingSavingID
		}

		err := sm.DeleteSaving(ctx, savingID, username)
		if err != nil {
			return err
		}

		return nil
	}
}

func setSavingGoalName(ctx context.Context, sgm SavingGoalManager, s *models.Saving) error {
	if s.SavingGoalID != nil && *s.SavingGoalID == "" {
		return nil
	}

	savingGoal, err := sgm.GetSavingGoal(ctx, s.Username, *s.SavingGoalID)
	if err != nil {
		return err
	}

	s.SavingGoalName = savingGoal.Name

	return nil
}

func setSavingGoalNames(ctx context.Context, sgm SavingGoalManager, username string, savings []*models.Saving) error {
	savingGoalsMap := make(map[string]*models.SavingGoal)

	savingGoals, err := sgm.GetSavingGoals(ctx, username)
	if err != nil {
		return err
	}

	for _, savingGoal := range savingGoals {
		savingGoalsMap[savingGoal.SavingGoalID] = savingGoal
	}

	for _, saving := range savings {
		if ignoreSaving(saving) {
			continue
		}

		savingGoal, ok := savingGoalsMap[*saving.SavingGoalID]
		if !ok {
			continue
		}

		saving.SavingGoalName = savingGoal.Name
	}

	return nil
}

func ignoreSaving(s *models.Saving) bool {
	return (s.SavingGoalID != nil && *s.SavingGoalID == "") || s.SavingGoalID == nil
}

func setSavingGoalNamesForSavingGoal(ctx context.Context, sgm SavingGoalManager, username, savingGoalID string, savings []*models.Saving) error {
	savingGoal, err := sgm.GetSavingGoal(ctx, username, savingGoalID)
	if err != nil {
		return err
	}

	for _, saving := range savings {
		saving.SavingGoalName = savingGoal.Name
	}

	return nil
}

func validateSavingPeriod(ctx context.Context, saving *models.Saving, username string, p PeriodManager) error {
	if saving.Period == nil {
		return nil
	}

	//TODO: get all periods
	periods, _, err := p.GetPeriods(ctx, username, "", 0)
	if err != nil {
		return err
	}

	for _, period := range periods {
		if period.ID == *saving.Period {
			return nil
		}
	}

	return models.ErrInvalidPeriod
}
