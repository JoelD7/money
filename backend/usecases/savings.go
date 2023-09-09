package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"math"
)

type SavingsManager interface {
	GetSavings(ctx context.Context, email string) ([]*models.Saving, error)
	CreateSaving(ctx context.Context, saving *models.Saving) error
	UpdateSaving(ctx context.Context, saving *models.Saving) error
}

func NewSavingsGetter(sm SavingsManager, l Logger) func(ctx context.Context, email string) ([]*models.Saving, error) {
	return func(ctx context.Context, email string) ([]*models.Saving, error) {
		err := validateEmail(email)
		if err != nil {
			l.Error("invalid_email_detected", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"s_email": email,
				}),
			})

			return nil, err
		}

		savings, err := sm.GetSavings(ctx, email)
		if err != nil {
			l.Error("savings_fetch_failed", err, []models.LoggerObject{
				l.MapToLoggerObject("user_data", map[string]interface{}{
					"s_email": email,
				}),
			})

			return nil, err
		}

		return savings, nil
	}
}

func NewSavingCreator(sm SavingsManager, l Logger) func(ctx context.Context, saving *models.Saving) error {
	return func(ctx context.Context, saving *models.Saving) error {
		err := validateSavingInput(saving)
		if err != nil {
			return fmt.Errorf("saving validation failed: %w", err)
		}

		err = sm.CreateSaving(ctx, saving)
		if err != nil {
			return fmt.Errorf("saving creation failed: %w", err)
		}

		return nil
	}
}

func NewSavingUpdater(sm SavingsManager) func(ctx context.Context, saving *models.Saving) error {
	return func(ctx context.Context, saving *models.Saving) error {
		err := validateSavingForUpdate(saving)
		if err != nil {
			return fmt.Errorf("saving validation failed: %w", err)
		}

		err = sm.UpdateSaving(ctx, saving)
		if err != nil {
			return err
		}

		return nil
	}
}

func validateSavingInput(saving *models.Saving) error {
	if *saving == (models.Saving{}) {
		return models.ErrInvalidRequestBody
	}

	err := validateEmail(saving.Email)
	if err != nil {
		return err
	}

	if saving.Amount <= 0 || saving.Amount > math.MaxFloat64 {
		return models.ErrInvalidAmount
	}

	return nil
}

func validateSavingForUpdate(saving *models.Saving) error {
	if err := validateSavingInput(saving); err != nil {
		return err
	}

	if saving.SavingID == "" {
		return models.ErrMissingSavingID
	}

	return nil
}
