package recurrent_expense_generator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-generator/handler"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/expenses"
	expenses_recurring "github.com/JoelD7/money/backend/storage/expenses-recurring"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

type testCase struct {
	username         string
	expectedExpenses int
	period           string
}

var (
	testCases = []*testCase{
		{
			username:         "e2e_test@gmail.com",
			expectedExpenses: 27,
			period:           "2021-09",
		},
		{
			username:         "e2e_test2@gmail.com",
			expectedExpenses: 11,
			period:           "2021-09",
		},
		{
			username:         "e2e_test3@gmail.com",
			expectedExpenses: 9,
			period:           "2021-09",
		},
		{
			username:         "e2e_test4@gmail.com",
			expectedExpenses: 15,
			period:           "2021-09",
		},
	}
)

func TestCron(t *testing.T) {
	var (
		expensesTableName          = env.GetString("EXPENSES_TABLE_NAME", "")
		expensesRecurringTableName = env.GetString("EXPENSES_RECURRING_TABLE_NAME", "")
		periodTableNameEnv         = env.GetString("PERIOD_TABLE_NAME", "")
		uniquePeriodTableNameEnv   = env.GetString("UNIQUE_PERIOD_TABLE_NAME", "")
		periodUserExpenseIndex     = env.GetString("PERIOD_USER_EXPENSE_INDEX", "")
	)

	c := require.New(t)

	ctx := context.Background()
	dynamoClient := setup.InitDynamoClient()
	log := logger.NewConsoleLogger("e2e-recurring-expense-generator")

	repo, err := expenses_recurring.NewExpenseRecurringDynamoRepository(dynamoClient, expensesRecurringTableName)
	c.Nil(err, "failed to create recurring expenses repository")

	periodRepo, err := period.NewDynamoRepository(dynamoClient, periodTableNameEnv, uniquePeriodTableNameEnv)
	c.Nil(err, "failed to create period repository")

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, expensesTableName, expensesRecurringTableName, periodUserExpenseIndex)
	c.Nil(err, "failed to create expenses repository")

	req := &handler.CronRequest{
		Log:          log,
		Repo:         repo,
		PeriodRepo:   periodRepo,
		ExpensesRepo: expensesRepo,
	}

	var expensesToDelete []*models.Expense
	var recExpenses []*models.ExpenseRecurring
	var periods []*models.Period

	t.Cleanup(func() {
		err := cleanup(req, &recExpenses, &expensesToDelete, &periods)
		c.Nil(err, "failed to cleanup test")
	})

	baseTime := time.Now()

	err = setupTestData(baseTime, req, &recExpenses, &periods)
	c.Nil(err, "failed to setupTestData test")
	c.NotEmpty(recExpenses, "recurring expenses array is empty")

	err = req.Process(ctx)
	c.Nil(err, "failed to process cron request")

	var userExpenses []*models.Expense

	for _, tc := range testCases {
		userExpenses, _, err = expensesRepo.GetExpenses(ctx, tc.username, "", 100)
		c.Nil(err, fmt.Sprintf("failed to get expenses for '%s'", tc.username))
		expensesToDelete = append(expensesToDelete, userExpenses...)
		c.Len(userExpenses, tc.expectedExpenses, fmt.Sprintf("expected %d expenses for '%s', got %d", tc.expectedExpenses, tc.username, len(userExpenses)))
	}
}

func setupTestData(baseTime time.Time, req *handler.CronRequest, recExpenses *[]*models.ExpenseRecurring, periods *[]*models.Period) error {
	re, err := setupRecurringExpenses(baseTime, req)
	if err != nil {
		return err
	}

	createdPeriods, err := createPeriods(req, baseTime)
	if err != nil {
		return err
	}

	*recExpenses = re
	*periods = createdPeriods

	return nil
}

func setupRecurringExpenses(baseTime time.Time, req *handler.CronRequest) ([]*models.ExpenseRecurring, error) {
	data, err := os.ReadFile("./samples/recurring_expenses.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read recurring_expenses.json: %v", err)
	}

	var recExpenses []*models.ExpenseRecurring
	err = json.Unmarshal(data, &recExpenses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sample JSON to recExpenses array: %v", err)
	}

	day := baseTime.Day()

	for _, expense := range recExpenses {
		//The cron uses the current day to generate the expenses
		expense.RecurringDay = day
		expense.CreatedDate = baseTime
	}

	err = req.Repo.BatchCreateExpenseRecurring(context.Background(), req.Log, recExpenses)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create recurring expenses: %v", err)
	}

	return recExpenses, nil
}

func createPeriods(req *handler.CronRequest, baseTime time.Time) ([]*models.Period, error) {
	var periodName string
	periods := make([]*models.Period, 0, len(testCases))

	for _, tc := range testCases {
		periodName = tc.period

		p := &models.Period{
			ID:        periodName,
			Username:  tc.username,
			Name:      &periodName,
			StartDate: baseTime,
			EndDate:   baseTime.AddDate(0, 1, 0),
		}

		_, err := req.PeriodRepo.CreatePeriod(context.Background(), p)
		if err != nil {
			return nil, fmt.Errorf("failed to create period '%s' for user '%s': %v", periodName, tc.username, err)
		}
		periods = append(periods, p)
	}

	return periods, nil
}

func cleanup(req *handler.CronRequest, recExpenses *[]*models.ExpenseRecurring, expensesToDelete *[]*models.Expense, periods *[]*models.Period) error {
	ctx := context.Background()
	var errs []error

	if recExpenses != nil && len(*recExpenses) > 0 {
		reErr := req.Repo.BatchDeleteExpenseRecurring(ctx, req.Log, *recExpenses)
		if reErr != nil {
			errs = append(errs, fmt.Errorf("failed to batch delete recurring expenses: %v", reErr))
		}
	}

	if expensesToDelete != nil && len(*expensesToDelete) > 0 {
		eErr := req.ExpensesRepo.BatchDeleteExpenses(ctx, *expensesToDelete)
		if eErr != nil {
			errs = append(errs, fmt.Errorf("failed to batch delete expenses: %v", eErr))
		}
	}

	if periods != nil && len(*periods) > 0 {
		pErr := req.PeriodRepo.BatchDeletePeriods(ctx, *periods)
		if pErr != nil {
			errs = append(errs, fmt.Errorf("failed to batch delete periods: %v", pErr))
		}
	}

	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	errString := ""
	for _, e := range errs {
		errString += e.Error() + " | "
	}

	return errors.New(errString)
}
