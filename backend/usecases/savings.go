package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"math"
	"math/rand"
	"time"
)

type SavingsManager interface {
	GetSavings(ctx context.Context, email string) ([]*models.Saving, error)
	CreateSaving(ctx context.Context, saving *models.Saving) error
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
			l.Error("saving_validation_failed", err, nil)

			return err
		}

		saving.SavingID = generateSavingID()
		saving.CreationDate = time.Now()

		err = sm.CreateSaving(ctx, saving)
		if err != nil {
			l.Error("create_saving_failed", err, nil)

			return err
		}

		return nil
	}
}

func validateSavingInput(saving *models.Saving) error {
	if *saving == (models.Saving{}) {
		return models.ErrEmptyRequestBody
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

func generateSavingID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return "SV" + string(b)
}
