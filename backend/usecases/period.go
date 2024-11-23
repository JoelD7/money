package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"math"
	"sync"
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
	DeletePeriod(ctx context.Context, periodID, username string) error
}

func NewPeriodCreator(pm PeriodManager, cache IncomePeriodCacheManager, log Logger) func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
	return func(ctx context.Context, username string, period *models.Period) (*models.Period, error) {
		if period.StartDate.After(period.EndDate) {
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

		err = cache.AddIncomePeriods(ctx, username, []string{periodID})
		if err != nil {
			log.Error("add_income_periods_failed", err, []models.LoggerObject{period})

			return nil, fmt.Errorf("couldn't add income periods to cache: %w", err)
		}

		err = sendPeriodToSQS(ctx, newPeriod)
		if err != nil {
			log.Error("send_period_to_sqs_failed", err, []models.LoggerObject{newPeriod})
		}

		return newPeriod, nil
	}
}

func sendPeriodToSQS(ctx context.Context, period *models.Period) error {
	missingExpensePeriodQueueURL := env.GetString("MISSING_EXPENSE_PERIOD_QUEUE_URL", "")

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

			income, err := im.GetAllIncomeByPeriod(ctx, username, periodID)
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

			expenses, err := em.GetAllExpensesByPeriod(ctx, username, periodID)
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
