package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type IncomeManager interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)
	GetIncomeByPeriod(ctx context.Context, username, periodID string) ([]*models.Income, error)
}

func NewIncomeCreator(im IncomeManager, pm PeriodManager) func(ctx context.Context, username string, income *models.Income) (*models.Income, error) {
	return func(ctx context.Context, username string, income *models.Income) (*models.Income, error) {
		err := validateIncomePeriod(ctx, username, income, pm)
		if err != nil {
			return nil, err
		}

		income.IncomeID = generateDynamoID("IN")
		income.Username = username
		income.CreatedDate = time.Now()

		newIncome, err := im.CreateIncome(ctx, income)
		if err != nil {
			return nil, err
		}

		return newIncome, nil
	}
}

func validateIncomePeriod(ctx context.Context, username string, income *models.Income, pm PeriodManager) error {
	if income.Period == nil {
		return nil
	}

	periods := make([]*models.Period, 0)
	curPeriods := make([]*models.Period, 0)
	nextKey := ""
	var err error

	for {
		curPeriods, nextKey, err = pm.GetPeriods(ctx, username, nextKey, 50)
		if err != nil {
			return fmt.Errorf("check if income period is valid failed: %v", err)
		}

		periods = append(periods, curPeriods...)

		if nextKey == "" {
			break
		}
	}

	for _, period := range periods {
		if period.ID == *income.Period {
			return nil
		}
	}

	return models.ErrInvalidPeriod
}
