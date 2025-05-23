package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"math"
	"sync"
	"time"
)

func NewPeriodCreator(pm PeriodManager, incomePeriodCache IncomePeriodCacheManager, resourceCache ResourceCacheManager, sgm SavingGoalManager, sm SavingsManager) func(ctx context.Context, username, idempotencyKey string, period *models.Period) (*models.Period, error) {
	return func(ctx context.Context, username, idempotencyKey string, period *models.Period) (*models.Period, error) {
		if period.StartDate.After(period.EndDate) {
			return nil, models.ErrStartDateShouldBeBeforeEndDate
		}

		period.Username = username
		period.CreatedDate = time.Now()

		newPeriod, err := CreateResource(ctx, resourceCache, idempotencyKey, func() (*models.Period, error) {
			newPeriod, err := pm.CreatePeriod(ctx, period)
			if err != nil {
				logger.Error("create_period_failed", err, models.Any("period", period))

				return nil, err
			}

			err = incomePeriodCache.AddIncomePeriods(ctx, username, []string{newPeriod.ID})
			if err != nil {
				logger.Error("add_income_periods_failed", err, models.Any("period", period))
			}

			err = sendPeriodToSQS(ctx, newPeriod)
			if err != nil {
				logger.Error("send_period_to_sqs_failed", err, models.Any("new_period", newPeriod))
			}

			err = generateRecurringSavings(ctx, username, newPeriod.Name, sgm, sm)
			if err != nil {
				logger.Error("generate_recurring_savings_failed", err, models.Any("new_period", newPeriod))
				return nil, err
			}

			return newPeriod, nil
		})

		if err != nil {
			return nil, err
		}

		return newPeriod, nil
	}
}

func generateRecurringSavings(ctx context.Context, username string, period *string, sgm SavingGoalManager, sm SavingsManager) error {
	goals, err := sgm.GetAllRecurringSavingGoals(ctx, username)
	if errors.Is(err, models.ErrSavingGoalsNotFound) {
		logger.Info("no_recurring_saving_goals_found", models.Any("username", username))
		return nil
	}

	if err != nil {
		return fmt.Errorf("couldn't get recurring saving goals: %w", err)
	}

	savingsToCreate := make([]*models.Saving, len(goals))

	for i, goal := range goals {
		savingsToCreate[i] = &models.Saving{
			Username:     username,
			Amount:       goal.RecurringAmount,
			CreatedDate:  time.Now(),
			SavingGoalID: &goal.SavingGoalID,
			Period:       period,
		}
	}

	return sm.BatchCreateSavings(ctx, savingsToCreate)
}

func sendPeriodToSQS(ctx context.Context, period *models.Period) error {
	missingExpensePeriodQueueURL := env.GetString("MISSING_EXPENSE_PERIOD_QUEUE_URL", "")
	isMissingExpensePeriodQueueEnabled := env.GetBool("ENABLE_MISSING_EXPENSE_PERIOD_QUEUE")

	if !isMissingExpensePeriodQueueEnabled {
		return nil
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("couldn't initialize SQS client: %w", err)
	}

	sqsClient := sqs.NewFromConfig(sdkConfig)

	msgBody, err := json.Marshal(&models.MissingExpensePeriodMessage{Period: period.ID, Username: period.Username})
	if err != nil {
		return fmt.Errorf("couldn't marshal message body: %w", err)
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(msgBody)),
		QueueUrl:    aws.String(missingExpensePeriodQueueURL),
	})

	if err != nil {
		return fmt.Errorf("couldn't send message to SQS: %w", err)
	}

	return nil
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

func NewPeriodDeleter(pm PeriodManager, cache IncomePeriodCacheManager) func(ctx context.Context, periodID, username string) error {
	return func(ctx context.Context, periodID, username string) error {
		err := cache.DeleteIncomePeriods(ctx, username, periodID)
		if err != nil {
			return fmt.Errorf("couldn't delete income periods from cache: %w", err)
		}

		return pm.DeletePeriod(ctx, periodID, username)
	}
}

func NewPeriodStatsGetter(em ExpenseManager, im IncomeRepository) func(ctx context.Context, username, periodID string) (*models.PeriodStat, error) {
	return func(ctx context.Context, username, periodID string) (*models.PeriodStat, error) {
		wg := sync.WaitGroup{}
		errChan := make(chan error, 2)
		totalIncome := 0.0
		categoryExpenseSummary := make([]*models.CategoryExpenseSummary, 0)

		wg.Add(1)
		go func() {
			defer func() { wg.Done() }()

			income, err := im.GetAllIncomeByPeriod(ctx, username, &models.QueryParameters{Period: periodID})
			if err != nil {
				errChan <- fmt.Errorf("couldn't get income for period: %w", err)
				return
			}

			for _, inc := range income {
				if inc.Amount != nil {
					totalIncome += *inc.Amount
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer func() { wg.Done() }()

			expenses, err := em.GetAllExpensesByPeriod(ctx, username, &models.QueryParameters{Period: periodID})
			if err != nil {
				errChan <- fmt.Errorf("couldn't get expenses for period: %w", err)
				return
			}

			categoryExpenses := make(map[string]float64)

			for _, expense := range expenses {
				if expense.CategoryID != nil && expense.Amount != nil {
					categoryExpenses[*expense.CategoryID] += *expense.Amount
				}
			}

			for category, amount := range categoryExpenses {
				categoryExpenseSummary = append(categoryExpenseSummary, &models.CategoryExpenseSummary{
					CategoryID: category,
					Total:      math.Round(amount*100) / 100,
				})
			}
		}()

		wg.Wait()
		close(errChan)

		select {
		case err := <-errChan:
			for e := range errChan {
				err = fmt.Errorf("%v: %w", err, e)
			}

			if err != nil {
				return nil, err
			}
		default:
		}

		return &models.PeriodStat{
			PeriodID:               periodID,
			TotalIncome:            totalIncome,
			CategoryExpenseSummary: categoryExpenseSummary,
		}, nil
	}
}
