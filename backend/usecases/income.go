package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

type IncomeRepository interface {
	CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error)

	GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error)
	GetAllIncome(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error)
	GetIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error)
	GetAllIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, error)
	GetAllIncomePeriods(ctx context.Context, username string) ([]string, error)
}

type IncomePeriodCacheManager interface {
	AddIncomePeriods(ctx context.Context, username string, periods []string) error
	GetIncomePeriods(ctx context.Context, username string) ([]string, error)
	DeleteIncomePeriods(ctx context.Context, username string, periods ...string) error
}

func NewIncomeCreator(im IncomeRepository, pm PeriodManager) func(ctx context.Context, username string, income *models.Income) (*models.Income, error) {
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

func NewIncomeGetter(im IncomeRepository) func(ctx context.Context, username, incomeID string) (*models.Income, error) {
	return func(ctx context.Context, username, incomeID string) (*models.Income, error) {
		return im.GetIncome(ctx, username, incomeID)
	}
}

func NewIncomeByPeriodGetter(repository IncomeRepository, cache IncomePeriodCacheManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, []string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, []string, error) {
		incomePeriods, err := getIncomePeriods(ctx, username, repository, cache)
		if err != nil {
			return nil, "", nil, err
		}

		income, nextKey, err := repository.GetIncomeByPeriod(ctx, username, params)
		if err != nil {
			return nil, "", nil, err
		}

		return income, nextKey, incomePeriods, nil
	}
}

func getIncomePeriods(ctx context.Context, username string, repository IncomeRepository, cache IncomePeriodCacheManager) ([]string, error) {
	incomePeriods, err := cache.GetIncomePeriods(ctx, username)
	if errors.Is(err, models.ErrIncomePeriodsNotFound) {
		incomePeriods, err = repository.GetAllIncomePeriods(ctx, username)
		if err != nil {
			return nil, err
		}

		err = cache.AddIncomePeriods(ctx, username, incomePeriods)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return incomePeriods, nil
}

func NewAllIncomeGetter(repository IncomeRepository, cache IncomePeriodCacheManager) func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, []string, error) {
	return func(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, []string, error) {
		incomePeriods, err := getIncomePeriods(ctx, username, repository, cache)
		if err != nil {
			return nil, "", nil, err
		}

		income, nextKey, err := repository.GetAllIncome(ctx, username, params)
		if err != nil {
			return nil, "", nil, err
		}

		return income, nextKey, incomePeriods, nil
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
