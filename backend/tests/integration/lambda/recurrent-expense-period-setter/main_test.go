package recurrent_expense_period_setter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/lambda/recurrent-expense-period-setter/handler"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	expensesRepo "github.com/JoelD7/money/backend/storage/expenses"
	periodRepo "github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var envConfig *models.EnvironmentConfiguration

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(err)
	}

	envConfig = env.GetEnvConfig()

	logger.InitLogger(logger.ConsoleImplementation)

	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	c := require.New(t)

	dynamoClient := setup.InitDynamoClient()
	ctx := context.Background()

	expensesRepository, err := expensesRepo.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating expenses repository failed")

	periodRepository, err := periodRepo.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating period repository failed")

	req := &handler.Request{
		ExpensesRepo: expensesRepository,
		PeriodRepo:   periodRepository,
	}

	username := "e2e_test@gmail.com"

	expensesList, err := loadExpenses()
	c.Nil(err, "loading expenses from file failed")
	c.NotEmpty(username, "username from loaded expenses is empty")

	startDate, err := time.Parse(time.DateOnly, "2023-09-01")
	c.Nil(err, "parsing start date failed")

	endDate, err := time.Parse(time.DateOnly, "2023-09-30")
	c.Nil(err, "parsing end date failed")

	period := &models.Period{
		ID:        "test-period",
		Name:      setup.StringPtr("test-period"),
		Username:  username,
		StartDate: startDate,
		EndDate:   endDate,
	}

	_, err = req.PeriodRepo.CreatePeriod(ctx, period)
	c.Nil(err, "creating period failed")

	msg := models.SQSMessage{
		SQSMessage: events.SQSMessage{
			Body: fmt.Sprintf(`{"period": "%s","username": "%s"}`, period.ID, username),
		},
	}

	err = req.ExpensesRepo.BatchCreateExpenses(ctx, expensesList)
	c.Nil(err, "batch creating expenses failed")

	t.Cleanup(func() {
		err = req.ExpensesRepo.BatchDeleteExpenses(ctx, expensesList)
		c.Nil(err, "batch deleting expenses failed")

		err = req.PeriodRepo.DeletePeriod(ctx, period.ID, period.Username)
		c.Nil(err, "deleting period failed")
	})

	err = req.ProcessMessage(ctx, msg)
	c.Nil(err)

	result, _, err := req.ExpensesRepo.GetExpensesByPeriod(ctx, period.Username, &models.QueryParameters{Period: period.ID, PageSize: 20})
	c.Nil(err, "getting expenses by period failed")
	c.Len(result, 18, fmt.Sprint("expected 18 expenses, got ", len(result)))
}

func loadExpenses() ([]*models.Expense, error) {
	data, err := os.ReadFile("./samples/expenses.json")
	if err != nil {
		return nil, err
	}

	var expensesList []*models.Expense
	err = json.Unmarshal(data, &expensesList)
	if err != nil {
		return nil, err
	}

	return expensesList, nil
}
