package expenses

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	repository "github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/tests/e2e/utils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetAllExpensesBetweenDates(t *testing.T) {
	var (
		expensesTableName          = env.GetString("EXPENSES_TABLE_NAME", "")
		expensesRecurringTableName = env.GetString("EXPENSES_RECURRING_TABLE_NAME", "")
	)

	c := require.New(t)

	dynamoClient := utils.InitDynamoClient()
	expensesRepo, err := repository.NewDynamoRepository(dynamoClient, expensesTableName, expensesRecurringTableName)

	expensesToCreate, err := loadExpenses()
	c.Nil(err, "failed to load expenses")

	ctx := context.Background()
	err = expensesRepo.BatchCreateExpenses(ctx, logger.NewConsoleLogger("test"), expensesToCreate)
	c.Nil(err, "failed to batch create expenses")

	t.Cleanup(func() {
		err = expensesRepo.BatchDeleteExpenses(ctx, expensesToCreate)
		c.Nil(err, "failed to batch delete expenses")
	})

	testCases := []struct {
		startDate     string
		endDate       string
		expectedCount int
	}{
		{
			startDate:     "2023-09-14",
			endDate:       "2023-09-15",
			expectedCount: 18,
		},
		{
			startDate:     "2023-10-14",
			endDate:       "2023-10-15",
			expectedCount: 12,
		},
		{
			startDate:     "2023-11-14",
			endDate:       "2023-11-15",
			expectedCount: 9,
		},
		{
			startDate:     "2023-12-14",
			endDate:       "2023-12-15",
			expectedCount: 19,
		},
	}

	var expenses []*models.Expense
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Expenses between %s and %s", tc.startDate, tc.endDate), func(t *testing.T) {
			c = require.New(t)

			expenses, err = expensesRepo.GetAllExpensesBetweenDates(ctx, "e2e_test@gmail.com", tc.startDate, tc.endDate)
			c.Nil(err, "failed to get all expenses between dates")
			c.Len(expenses, tc.expectedCount, fmt.Sprintf("expected %d expenses, got %d", tc.expectedCount, len(expenses)))
		})

	}
}

func loadExpenses() ([]*models.Expense, error) {
	data, err := os.ReadFile("./samples/expenses.json")
	if err != nil {
		return nil, err
	}

	var expenses []*models.Expense
	err = json.Unmarshal(data, &expenses)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
