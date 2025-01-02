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
	run()
}

func run() {
	log := logger.NewConsoleLogger("expenses-test-data")
	ctx := context.Background()

	dynamoClient := dynamo.InitClient(ctx)
	envConfig, err := env.LoadEnv(ctx)
	if err != nil {
		panic(fmt.Errorf("loading environment variables failed: %v", err))
	}

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, envConfig)
	if err != nil {
		panic(fmt.Errorf("creating expenses repository failed: %v", err))
	}

	data, err := os.ReadFile("./samples/expenses.json")
	if err != nil {
		panic(fmt.Errorf("reading expenses file failed: %w", err))
	}

	var expensesList []*models.Expense
	err = json.Unmarshal(data, &expensesList)
	if err != nil {
		panic(fmt.Errorf("unmarshalling expenses list failed: %w", err))
	}

	err = expensesRepo.BatchCreateExpenses(ctx, log, expensesList)
	if err != nil {
		panic(fmt.Errorf("batch creating expenses failed: %w", err))
	}
}
