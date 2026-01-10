package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"time"
)

func NewIncomeCreator(im IncomeRepository, pm PeriodManager, cache ResourceCacheManager) func(ctx context.Context, username, idempotencyKey string, income *models.Income) (*models.Income, error) {
	return func(ctx context.Context, username, idempotencyKey string, income *models.Income) (*models.Income, error) {
		return CreateResource(ctx, cache, idempotencyKey, func() (*models.Income, error) {
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
		})
	}
}

func NewIncomeGetter(im IncomeRepository, pm PeriodManager) func(ctx context.Context, username, incomeID string) (*models.Income, error) {
	return func(ctx context.Context, username, incomeID string) (*models.Income, error) {
		income, err := im.GetIncome(ctx, username, incomeID)
		if err != nil {
			return nil, err
		}

		err = setEntitiesPeriods(ctx, pm, income)
		if err != nil {
			return nil, fmt.Errorf("couldn't set periods for income: %w", err)
		}

		return income, nil
	}
}

func NewIncomeByPeriodGetter(repository IncomeRepository, cache IncomePeriodCacheManager, pm PeriodManager) func(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, []string, error) {
	return func(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, []string, error) {
		incomePeriods, err := getIncomePeriods(ctx, username, repository, cache)
		if err != nil {
			return nil, "", nil, err
		}

		income, nextKey, err := repository.GetIncomeByPeriod(ctx, username, params)
		if err != nil {
			return nil, "", nil, err
		}

		periodManipulator := make([]PeriodHolder, len(income))
		for i := 0; i < len(income); i++ {
			periodManipulator[i] = income[i]
		}

		err = setEntitiesPeriods(ctx, pm, periodManipulator...)
		if err != nil {
			return nil, "", nil, fmt.Errorf("couldn't set periods for income: %w", err)
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

func NewAllIncomeGetter(repository IncomeRepository, cache IncomePeriodCacheManager, pm PeriodManager) func(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, []string, error) {
	return func(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, []string, error) {
		incomePeriods, err := getIncomePeriods(ctx, username, repository, cache)
		if err != nil {
			return nil, "", nil, err
		}

		income, nextKey, err := repository.GetAllIncome(ctx, username, params)
		if err != nil {
			return nil, "", nil, err
		}

		periodManipulator := make([]PeriodHolder, len(income))
		for i := 0; i < len(income); i++ {
			periodManipulator[i] = income[i]
		}

		err = setEntitiesPeriods(ctx, pm, periodManipulator...)
		if err != nil {
			return nil, "", nil, fmt.Errorf("couldn't set periods for income: %w", err)
		}

		return income, nextKey, incomePeriods, nil
	}
}

func validateIncomePeriod(ctx context.Context, username string, income *models.Income, pm PeriodManager) error {
	if income == nil || income.PeriodID == nil {
		return nil
	}

	var err error

	_, err = pm.GetPeriod(ctx, username, *income.PeriodID)
	if errors.Is(err, models.ErrPeriodNotFound) {
		return models.ErrInvalidPeriod
	}

	if err != nil {
		return fmt.Errorf("check if expense period is valid failed: %v", err)
	}

	return nil
}
