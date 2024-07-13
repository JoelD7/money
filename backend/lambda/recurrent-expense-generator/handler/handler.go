package handler

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	expenses_recurring "github.com/JoelD7/money/backend/storage/expenses-recurring"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/shared"
	"github.com/JoelD7/money/backend/usecases"
	"sync"
	"time"
)

var (
	once sync.Once

	expensesTableName          = env.GetString("EXPENSES_TABLE_NAME", "")
	expensesRecurringTableName = env.GetString("EXPENSES_RECURRING_TABLE_NAME", "")
)

type CronRequest struct {
	Log          logger.LogAPI
	Repo         expenses_recurring.Repository
	PeriodRepo   period.Repository
	ExpensesRepo expenses.Repository

	err          error
	startingTime time.Time
}

func Handle(ctx context.Context) error {
	req := new(CronRequest)
	var err error

	stackTrace, ctxError := shared.ExecuteLambda(ctx, func(ctx context.Context) {
		err = req.init(ctx)
		if err != nil {
			return
		}

		defer req.finish()

		err = req.Process(ctx)
	})

	if ctxError != nil {
		req.Log.Error("request_timeout", ctxError, []models.LoggerObject{
			req.Log.MapToLoggerObject("stack", map[string]interface{}{
				"s_trace": stackTrace,
			}),
		})
	}

	if err != nil {
		req.Log.Error("request_error", err, nil)

		return err
	}

	return nil
}

func (req *CronRequest) init(ctx context.Context) error {
	var err error
	once.Do(func() {
		req.Log = logger.NewLogger()

		dynamoClient := dynamo.InitClient(ctx)
		req.Repo, err = expenses_recurring.NewExpenseRecurringDynamoRepository(dynamoClient, expensesRecurringTableName)
		if err != nil {
			return
		}

		req.PeriodRepo = period.NewDynamoRepository(dynamoClient)

		req.ExpensesRepo, err = expenses.NewDynamoRepository(dynamoClient, expensesTableName, expensesRecurringTableName)
		if err != nil {
			return
		}
	})
	req.startingTime = time.Now()

	return err
}

func (req *CronRequest) finish() {
	defer func() {
		err := req.Log.Finish()
		if err != nil {
			panic(err)
		}
	}()

	req.Log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *CronRequest) Process(ctx context.Context) error {
	day := time.Now().Day()

	recExpenses, err := req.Repo.ScanExpensesForDay(ctx, day)
	if err != nil {
		req.err = err
		req.Log.Error("scan_expenses_for_day_failed", err, []models.LoggerObject{
			req.Log.MapToLoggerObject("run_information", map[string]interface{}{
				"i_day": day,
			}),
		})

		return err
	}

	recExpensesByUser := make(map[string][]*models.ExpenseRecurring)
	for _, expense := range recExpenses {
		recExpensesByUser[expense.Username] = append(recExpensesByUser[expense.Username], expense)
	}

	for username, userRecurringExpenses := range recExpensesByUser {
		err = req.createExpenses(ctx, username, userRecurringExpenses)
		if err != nil {
			req.Log.Error("create_expenses_failed", err, []models.LoggerObject{
				req.Log.MapToLoggerObject("run_information", map[string]interface{}{
					"s_username": username,
				}),
			})
		}
	}

	return nil
}

func (req *CronRequest) createExpenses(ctx context.Context, username string, recExpenses []*models.ExpenseRecurring) error {
	lastPeriod, err := req.PeriodRepo.GetLastPeriod(ctx, username)
	if err != nil {
		return fmt.Errorf("get last period failed: %v", err)
	}

	expensesToCreate := make([]*models.Expense, 0, len(recExpenses))

	for _, expense := range recExpenses {
		if isRecurringExpenseWithinPeriod(expense, lastPeriod) {
			expensesToCreate = append(expensesToCreate, &models.Expense{
				Username:   expense.Username,
				CategoryID: expense.CategoryID,
				Amount:     &expense.Amount,
				Name:       &expense.Name,
				Notes:      expense.Notes,
				Period:     *lastPeriod.Name,
			})
		}
	}

	batchCreateExpenses := usecases.NewBatchExpensesCreator(req.ExpensesRepo, req.Log)
	err = batchCreateExpenses(ctx, expensesToCreate)
	if err != nil {
		return fmt.Errorf("batch create expenses failed: %v", err)
	}

	return nil
}

func isRecurringExpenseWithinPeriod(re *models.ExpenseRecurring, p *models.Period) bool {
	return (re.CreatedDate.After(p.StartDate) || re.CreatedDate.Equal(p.StartDate)) && (re.CreatedDate.Before(p.EndDate) || re.CreatedDate.Equal(p.EndDate))
}
