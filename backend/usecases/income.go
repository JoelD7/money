package usecases

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type IncomeManager interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)
	GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error)
	GetAllIncome(ctx context.Context, username, startKey string, pageSize int) ([]*models.Income, string, error)
	GetIncomeByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Income, string, error)
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

func NewIncomeGetter(im IncomeManager) func(ctx context.Context, username, incomeID string) (*models.Income, error) {
	return func(ctx context.Context, username, incomeID string) (*models.Income, error) {
		return im.GetIncome(ctx, username, incomeID)
	}
}

func NewIncomeByPeriodGetter(im IncomeManager) func(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Income, string, error) {
	return func(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Income, string, error) {
		return im.GetIncomeByPeriod(ctx, username, periodID, startKey, pageSize)
	}
}

func NewAllIncomeGetter(im IncomeManager) func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Income, string, error) {
	return func(ctx context.Context, username, startKey string, pageSize int) ([]*models.Income, string, error) {
		return im.GetAllIncome(ctx, username, startKey, pageSize)
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
