package expenses_recurring

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var awsRegion = env.GetString("REGION", "us-east-1")

func initDynamoClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func TestScanExpensesForDay(t *testing.T) {
	//c := require.New(t)
	//
	//dynamoClient := initDynamoClient()
	//
	//repository := NewExpenseRecurringDynamoRepository(dynamoClient)
	//
	//createExpenses(c, repository)
	//
	//var expenses []*models.ExpenseRecurring
	//var err error
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 15)
	//c.Nil(err)
	//c.Len(expenses, 2)
	//c.False(areRepeated(expenses))
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 5)
	//c.Nil(err)
	//c.Len(expenses, 1)
	//c.False(areRepeated(expenses))
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 1)
	//c.Nil(err)
	//c.Len(expenses, 11)
	//c.False(areRepeated(expenses))
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 10)
	//c.Nil(err)
	//c.Len(expenses, 1)
	//c.False(areRepeated(expenses))
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 20)
	//c.Nil(err)
	//c.Len(expenses, 2)
	//c.False(areRepeated(expenses))
	//
	//expenses, err = repository.ScanExpensesForDay(context.Background(), 25)
	//c.Nil(err)
	//c.Len(expenses, 1)
	//c.False(areRepeated(expenses))
}

func createExpenses(c *require.Assertions, repository *DynamoRepository) {
	b := []byte(`[
	 {
	   "id": "gym membership",
	   "username": "test@gmail.com",
	   "category_id": "CTGcSuhjzVmu3WrHLKD5fhS",
	   "amount": 50.0,
	   "recurring_day": 1,
	   "name": "Gym Membership",
	   "notes": "Monthly gym fee"
	 },
	 {
	   "id": "netflix subscription",
	   "username": "test@gmail.com",
	   "category_id": "CTGGyouAaIPPWKzxpyxHACS",
	   "amount": 15.0,
	   "recurring_day": 1,
	   "name": "Netflix Subscription",
	   "notes": "Monthly streaming service"
	 },
	 {
	   "id": "spotify premium",
	   "username": "test@gmail.com",
	   "category_id": "CTGGyouAaIPPWKzxpyxHACS",
	   "amount": 10.0,
	   "recurring_day": 1,
	   "name": "Spotify Premium",
	   "notes": "Monthly music streaming"
	 },
	 {
	   "id": "yoga class",
	   "username": "test@gmail.com",
	   "category_id": "CTGcSuhjzVmu3WrHLKD5fhS",
	   "amount": 20.0,
	   "recurring_day": 5,
	   "name": "Yoga Class",
	   "notes": "Weekly yoga sessions"
	 },
	 {
	   "id": "amazon prime",
	   "username": "test@gmail.com",
	   "category_id": "CTGGyouAaIPPWKzxpyxHACS",
	   "amount": 12.99,
	   "recurring_day": 1,
	   "name": "Amazon Prime",
	   "notes": "Monthly subscription"
	 },
	 {
	   "id": "internet bill",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 60.0,
	   "recurring_day": 10,
	   "name": "Internet Bill",
	   "notes": "Monthly internet service"
	 },
	 {
	   "id": "electricity bill",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 100.0,
	   "recurring_day": 15,
	   "name": "Electricity Bill",
	   "notes": "Monthly electricity charges"
	 },
	 {
	   "id": "water bill",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 30.0,
	   "recurring_day": 20,
	   "name": "Water Bill",
	   "notes": "Monthly water charges"
	 },
	 {
	   "id": "phone bill",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 40.0,
	   "recurring_day": 25,
	   "name": "Phone Bill",
	   "notes": "Monthly phone charges"
	 },
	 {
	   "id": "car insurance",
	   "username": "test@gmail.com",
	   "category_id": "CTG2Tb6hKgnr2mva6hwpviA",
	   "amount": 120.0,
	   "recurring_day": 1,
	   "name": "Car Insurance",
	   "notes": "Monthly car insurance"
	 },
	 {
	   "id": "health insurance",
	   "username": "test@gmail.com",
	   "category_id": "CTG2Tb6hKgnr2mva6hwpviA",
	   "amount": 200.0,
	   "recurring_day": 1,
	   "name": "Health Insurance",
	   "notes": "Monthly CTGcSuhjzVmu3WrHLKD5fhS insurance"
	 },
	 {
	   "id": "rent",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 1500.0,
	   "recurring_day": 1,
	   "name": "Rent",
	   "notes": "Monthly house rent"
	 },
	 {
	   "id": "mortgage",
	   "username": "test@gmail.com",
	   "category_id": "CTGtOZzZ2oOxVT9ahBg8WLw",
	   "amount": 2000.0,
	   "recurring_day": 1,
	   "name": "Mortgage",
	   "notes": "Monthly house mortgage"
	 },
	 {
	   "id": "car loan",
	   "username": "test@gmail.com",
	   "category_id": "CTGQ5nsMlz8CmBDX4lTSJDx",
	   "amount": 300.0,
	   "recurring_day": 15,
	   "name": "Car Loan",
	   "notes": "Monthly car loan payment"
	 },
	 {
	   "id": "student loan",
	   "username": "test@gmail.com",
	   "category_id": "CTGQ5nsMlz8CmBDX4lTSJDx",
	   "amount": 400.0,
	   "recurring_day": 1,
	   "name": "Student Loan",
	   "notes": "Monthly student loan payment"
	 },
	 {
	   "id": "credit card",
	   "username": "test@gmail.com",
	   "category_id": "CTGQ5nsMlz8CmBDX4lTSJDx",
	   "amount": 500.0,
	   "recurring_day": 20,
	   "name": "Credit Card",
	   "notes": "Monthly credit card payment"
	 },
	 {
	   "id": "netflix subscription",
	   "username": "test@gmail.com",
	   "category_id": "CTGGyouAaIPPWKzxpyxHACS",
	   "amount": 15.0,
	   "recurring_day": 1,
	   "name": "Netflix Subscription",
	   "notes": "Monthly streaming service"
	 },
	 {
	   "id": "gym membership",
	   "username": "test@gmail.com",
	   "category_id": "CTGcSuhjzVmu3WrHLKD5fhS",
	   "amount": 50.0,
	   "recurring_day": 1,
	   "name": "Gym Membership",
	   "notes": "Monthly gym fee"
	 },
	 {
	   "id": "insurance premium",
	   "username": "test@gmail.com",
	   "category_id": "CTG2Tb6hKgnr2mva6hwpviA",
	   "amount": 100.0,
	   "recurring_day": 1,
	   "name": "Insurance Premium",
	   "notes": "Monthly insurance premium"
	 },
	 {
	   "id": "software subscription",
	   "username": "test@gmail.com",
	   "category_id": "CTGGyouAaIPPWKzxpyxHACS",
	   "amount": 25.0,
	   "recurring_day": 1,
	   "name": "Software Subscription",
	   "notes": "Monthly software subscription"
	 }
	]
	`)

	expenses := make([]*models.ExpenseRecurring, 0)

	err := json.Unmarshal(b, &expenses)
	c.Nil(err)

	c.Len(expenses, 20)

	for _, expense := range expenses {
		expense.CreatedDate = time.Now()
		_, err = repository.CreateExpenseRecurring(context.Background(), expense)
		c.Nil(err)
	}
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
