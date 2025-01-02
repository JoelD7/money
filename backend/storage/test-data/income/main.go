package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"os"
)

func main() {
	// Keep this call commented to avoid accidental executions
	//run()
}

func run() {
	ctx := context.Background()

	dynamoClient := dynamo.InitClient(ctx)
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment variables failed: %v", err))
	}

	envConfig := &models.EnvironmentConfiguration{
		IncomeTable:                 env.GetString("INCOME_TABLE_NAME", ""),
		PeriodUserIncomeIndex:       env.GetString("PERIOD_USER_INCOME_INDEX", ""),
		PeriodUserNameIncomeIDIndex: env.GetString("PERIOD_USER_NAME_INCOME_ID_INDEX", ""),
		PeriodUserAmountIndex:       env.GetString("PERIOD_USER_AMOUNT_INDEX", ""),
		PeriodUserCreatedDateIndex:  env.GetString("PERIOD_USER_CREATED_DATE_INDEX", ""),
		UsernameCreatedDateIndex:    env.GetString("USERNAME_CREATED_DATE_INDEX", ""),
	}

	incomeRepo, err := income.NewDynamoRepository(dynamoClient, envConfig)
	if err != nil {
		panic(fmt.Errorf("creating income repository failed: %v", err))
	}

	data, err := os.ReadFile("/Users/joelfabian/go/src/github.com/JoelD7/money/backend/storage/test-data/income/income.json")
	if err != nil {
		panic(fmt.Errorf("reading income file failed: %w", err))
	}

	var incomelist []*models.Income
	err = json.Unmarshal(data, &incomelist)
	if err != nil {
		panic(fmt.Errorf("unmarshalling income list failed: %w", err))
	}

	err = incomeRepo.BatchCreateIncome(ctx, incomelist)
	if err != nil {
		panic(fmt.Errorf("batch creating nicome failed: %w", err))
	}

	fmt.Println("Income data created successfully")
}
