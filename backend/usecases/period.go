package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"strconv"
	"strings"
	"time"
)

const (
	// yearlyPeriods is the number of periods in a year.
	yearlyPeriods = 12
)

type PeriodManager interface {
	CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error)
	GetPeriod(ctx context.Context, username, period string) (*models.Period, error)
	GetLastPeriod(ctx context.Context, username string) (*models.Period, error)
	GetPeriods(ctx context.Context, username string) ([]*models.Period, error)
}

func NewPeriodCreator(pm PeriodManager) func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
	return func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
		err := validatePeriod(period)
		if err != nil {
			return nil, err
		}

		periodID, err := generateNewPeriodID(ctx, pm, username)
		if err != nil {
			return nil, fmt.Errorf("generate new period ID failed: %w", err)
		}

		period.ID = periodID
		period.Username = username
		period.CreatedDate = time.Now()

		return pm.CreatePeriod(ctx, period)
	}
}

func generateNewPeriodID(ctx context.Context, pm PeriodManager, username string) (string, error) {
	lastPeriod, err := pm.GetLastPeriod(ctx, username)
	if err != nil && !errors.Is(err, models.ErrPeriodsNotFound) {
		return "", err
	}

	if errors.Is(err, models.ErrPeriodsNotFound) {
		return fmt.Sprintf("%s-%s", strconv.Itoa(time.Now().Year()), "1"), nil
	}

	errMalformedPeriodID := fmt.Errorf("malformed period ID: %s", lastPeriod.ID)

	periodParts := strings.Split(lastPeriod.ID, "-")
	if len(periodParts) != 2 {
		return "", errMalformedPeriodID
	}

	periodNumber, err := strconv.Atoi(periodParts[1])
	if err != nil {
		return "", fmt.Errorf("%v: %v", lastPeriod.ID, err)
	}

	periodYear, err := strconv.Atoi(periodParts[0])
	if err != nil {
		return "", fmt.Errorf("%v: %v", lastPeriod.ID, err)
	}

	if periodNumber+1 > yearlyPeriods {
		periodNumber = 1
		periodYear++
	} else {
		periodNumber++
	}

	return fmt.Sprintf("%s-%s", strconv.Itoa(periodYear), strconv.Itoa(periodNumber)), nil
}

func validatePeriod(period *models.Period) error {
	if period.StartDate.After(period.EndDate.Time) {
		return models.ErrStartDateShouldBeBeforeEndDate
	}

	return nil
}
