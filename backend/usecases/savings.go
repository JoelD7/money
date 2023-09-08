package usecases

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type SavingsManager interface {
	GetSavings(ctx context.Context, email string) ([]*models.Saving, error)
}

func NewSavingsGetter(sm SavingsManager, l Logger) func(ctx context.Context, email string) ([]*models.Saving, error) {
	return func(ctx context.Context, email string) ([]*models.Saving, error) {
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
