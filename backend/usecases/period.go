package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

const (
	periodPrefix = "PRD"
)

type PeriodManager interface {
	CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error)
	UpdatePeriod(ctx context.Context, period *models.Period) error
	GetPeriod(ctx context.Context, username, period string) (*models.Period, error)
	GetLastPeriod(ctx context.Context, username string) (*models.Period, error)
	GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error)
}

func NewPeriodCreator(pm PeriodManager, log Logger) func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
	return func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
		if period.StartDate.After(period.EndDate.Time) {
			return nil, models.ErrStartDateShouldBeBeforeEndDate
		}

		periodID := generateDynamoID(periodPrefix)

		period.ID = periodID
		period.Username = username
		period.CreatedDate = time.Now()

		newPeriod, err := pm.CreatePeriod(ctx, period)
		if err != nil {
			log.Error("create_period_failed", err, []models.LoggerObject{period})

			return nil, err
		}

		return newPeriod, nil
	}
}

func NewPeriodUpdater(pm PeriodManager) func(ctx context.Context, username, periodID string, period *models.Period) (*models.Period, error) {
	return func(ctx context.Context, username, periodID string, period *models.Period) (*models.Period, error) {
		err := pm.UpdatePeriod(ctx, period)
		if err != nil {
			return nil, err
		}

		updatedPeriod, err := pm.GetPeriod(ctx, username, periodID)
		if err != nil {
			return nil, fmt.Errorf("get updated period failed: %w", err)
		}

		return updatedPeriod, nil
	}
}

func NewPeriodGetter(pm PeriodManager) func(ctx context.Context, username, periodID string) (*models.Period, error) {
	return func(ctx context.Context, username, periodID string) (*models.Period, error) {
		return pm.GetPeriod(ctx, username, periodID)
	}
}

func NewPeriodsGetter(pm PeriodManager) func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error) {
	return func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error) {
		return pm.GetPeriods(ctx, username, startKey, pageSize)
	}
}
