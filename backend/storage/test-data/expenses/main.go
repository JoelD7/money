package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"os"
)

func main() {
	// Keep this call commented to avoid accidental executions
	//run()
}

func run() {
	logger.InitLogger(logger.ConsoleImplementation)

	ctx := context.Background()

	dynamoClient := dynamo.InitClient(ctx)
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment variables failed: %v", err))
	}

	envConfig := &models.EnvironmentConfiguration{
		ExpensesTable:                env.GetString("EXPENSES_TABLE_NAME", ""),
		ExpensesRecurringTable:       env.GetString("EXPENSES_RECURRING_TABLE_NAME", ""),
		PeriodUserExpenseIndex:       env.GetString("PERIOD_USER_EXPENSE_INDEX", ""),
		UsersTable:                   env.GetString("USERS_TABLE_NAME", ""),
		PeriodUserCreatedDateIndex:   env.GetString("PERIOD_USER_CREATED_DATE_INDEX", ""),
		UsernameCreatedDateIndex:     env.GetString("USERNAME_CREATED_DATE_INDEX", ""),
		PeriodUserNameExpenseIDIndex: env.GetString("PERIOD_USER_NAME_EXPENSE_ID_INDEX", ""),
		PeriodUserAmountIndex:        env.GetString("PERIOD_USER_AMOUNT_INDEX", ""),
	}

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, envConfig)
	if err != nil {
		panic(fmt.Errorf("creating expenses repository failed: %v", err))
	}

	data, err := os.ReadFile("/Users/joelfabian/go/src/github.com/JoelD7/money/backend/storage/test-data/expenses/expenses.json")
	if err != nil {
		panic(fmt.Errorf("reading expenses file failed: %w", err))
	}

	var expensesList []*models.Expense
	err = json.Unmarshal(data, &expensesList)
	if err != nil {
		panic(fmt.Errorf("unmarshalling expenses list failed: %w", err))
	}

	err = expensesRepo.BatchCreateExpenses(ctx, expensesList)
	if err != nil {
		panic(fmt.Errorf("batch creating expenses failed: %w", err))
	}
}
