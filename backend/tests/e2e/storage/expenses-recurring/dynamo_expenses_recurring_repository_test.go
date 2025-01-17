package expenses_recurring

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	repo "github.com/JoelD7/money/backend/storage/expenses-recurring"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestScanExpensesForDay(t *testing.T) {
	c := require.New(t)

	logger.InitLogger(logger.ConsoleImplementation)

	var err error

	dynamoClient := setup.InitDynamoClient()
	expensesRecurringTableName := env.GetString("EXPENSES_RECURRING_TABLE_NAME", "")

	repository, err := repo.NewExpenseRecurringDynamoRepository(dynamoClient, expensesRecurringTableName)

	createExpenses(c, repository)

	var expenses []*models.ExpenseRecurring

	expenses, err = repository.ScanExpensesForDay(context.Background(), 15)
	c.Nil(err)
	c.Len(expenses, 11)
	c.False(areRepeated(expenses))

	expenses, err = repository.ScanExpensesForDay(context.Background(), 10)
	c.Nil(err)
	c.Len(expenses, 6)
	c.False(areRepeated(expenses))

	expenses, err = repository.ScanExpensesForDay(context.Background(), 19)
	c.Nil(err)
	c.Len(expenses, 8)
	c.False(areRepeated(expenses))

	expenses, err = repository.ScanExpensesForDay(context.Background(), 1)
	c.Nil(err)
	c.Len(expenses, 2)
	c.False(areRepeated(expenses))
}

func createExpenses(c *require.Assertions, repository *repo.DynamoRepository) {
	data, err := os.ReadFile("./samples/recurring_expenses.json")
	c.Nil(err, fmt.Sprintf("failed to read recurring_expenses.json: %v", err))

	var recExpenses []*models.ExpenseRecurring
	err = json.Unmarshal(data, &recExpenses)
	c.Nil(err, fmt.Sprintf("failed to unmarshal sample JSON to recExpenses array: %v", err))
	c.Len(recExpenses, 27)

	expensesPerDay := map[int]int{}

	for _, expense := range recExpenses {
		if _, ok := expensesPerDay[expense.RecurringDay]; !ok {
			expensesPerDay[expense.RecurringDay] = 1
			continue
		}
		expensesPerDay[expense.RecurringDay]++
	}

	for _, expense := range recExpenses {
		expense.CreatedDate = time.Now()
	}

	err = repository.BatchCreateExpenseRecurring(context.Background(), recExpenses)
	c.Nil(err)
}

func areRepeated(expenses []*models.ExpenseRecurring) bool {
	seen := make(map[string]bool)
	for _, expense := range expenses {
		if seen[expense.ID] {
			return true
		}
		seen[expense.ID] = true
	}
	return false
}
