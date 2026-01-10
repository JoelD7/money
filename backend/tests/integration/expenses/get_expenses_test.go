package expenses

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

var (
	envConfig *models.EnvironmentConfiguration
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(err)
	}

	logger.InitLogger(logger.ConsoleImplementation)

	envConfig = env.GetEnvConfig()

	os.Exit(m.Run())
}

func TestGetByPeriod(t *testing.T) {
	c := require.New(t)
	ctx := context.Background()

	expensesRepo, err := expenses.NewDynamoRepository(dynamo.InitClient(ctx), envConfig)
	c.NoError(err)

	usersRepo, err := users.NewDynamoRepository(dynamo.InitClient(ctx), envConfig.UsersTable)
	c.NoError(err)

	username := "e2e_test@gmail.com"
	periodID := "2021-09"

	ger := &handlers.GetExpensesRequest{
		ExpensesRepo: expensesRepo,
		Username:     username,
		UserRepo:     usersRepo,
	}

	testUser := &models.User{
		Username:      username,
		CurrentPeriod: periodID,
	}

	err = setup.CreateUser(ctx, usersRepo, testUser, t)
	c.Nil(err, "creating user failed")

	_, err = setup.CreateExpensesEntries(ctx, expensesRepo, "", t)
	c.Nil(err, "creating expenses failed")

	apigwReq := new(apigateway.Request)

	t.Run("Sorted by created_date - ASC", func(t *testing.T) {
		ger.ExpenseQueryParameters = &models.ExpenseQueryParameters{
			SortBy: "created_date",
			Period: periodID,
		}

		res, err := ger.GetByPeriod(ctx, apigwReq)
		c.NoError(err)
		c.Equal(res.StatusCode, http.StatusOK)

		var expensesResponse handlers.ExpensesResponse
		err = json.Unmarshal([]byte(res.Body), &expensesResponse)
		c.NoError(err)

		expected := []models.Expense{
			{
				ExpenseID:   "insurance premium",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(100),
				Name:        stringPtr("Insurance Premium"),
				Notes:       "Monthly insurance premium",
				CreatedDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gym membership",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(50),
				Name:        stringPtr("Gym Membership"),
				Notes:       "Monthly gym membership",
				CreatedDate: time.Date(2021, 9, 2, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "netflix subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(12.99),
				Name:        stringPtr("Netflix Subscription"),
				Notes:       "Monthly Netflix subscription",
				CreatedDate: time.Date(2021, 9, 3, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "spotify subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(9.99),
				Name:        stringPtr("Spotify Subscription"),
				Notes:       "Monthly Spotify subscription",
				CreatedDate: time.Date(2021, 9, 4, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "car insurance",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(75),
				Name:        stringPtr("Car Insurance"),
				Notes:       "Monthly car insurance premium",
				CreatedDate: time.Date(2021, 9, 5, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "electricity bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(60),
				Name:        stringPtr("Electricity Bill"),
				Notes:       "Monthly electricity bill",
				CreatedDate: time.Date(2021, 9, 6, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "internet bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(40),
				Name:        stringPtr("Internet Bill"),
				Notes:       "Monthly internet bill",
				CreatedDate: time.Date(2021, 9, 7, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "water bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(30),
				Name:        stringPtr("Water Bill"),
				Notes:       "Monthly water bill",
				CreatedDate: time.Date(2021, 9, 8, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "phone bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(25),
				Name:        stringPtr("Phone Bill"),
				Notes:       "Monthly phone bill",
				CreatedDate: time.Date(2021, 9, 9, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "credit card payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(200),
				Name:        stringPtr("Credit Card Payment"),
				Notes:       "Monthly credit card payment",
				CreatedDate: time.Date(2021, 9, 10, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "student loan",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(150),
				Name:        stringPtr("Student Loan"),
				Notes:       "Monthly student loan payment",
				CreatedDate: time.Date(2021, 9, 11, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "rent payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(800),
				Name:        stringPtr("Rent Payment"),
				Notes:       "Monthly rent payment",
				CreatedDate: time.Date(2021, 9, 12, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gas bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(20),
				Name:        stringPtr("Gas Bill"),
				Notes:       "Monthly gas bill",
				CreatedDate: time.Date(2021, 9, 13, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
		}

		for i, expense := range expensesResponse.Expenses {
			c.Equal(expected[i].ExpenseID, expense.ExpenseID)
		}
	})

	t.Run("Sorted by created_date - DESC", func(t *testing.T) {
		ger.ExpenseQueryParameters = &models.ExpenseQueryParameters{
			SortBy:   "created_date",
			SortType: "desc",
			Period:   periodID,
		}

		res, err := ger.GetByPeriod(ctx, apigwReq)
		c.NoError(err)
		c.Equal(res.StatusCode, http.StatusOK)

		var expensesResponse handlers.ExpensesResponse
		err = json.Unmarshal([]byte(res.Body), &expensesResponse)
		c.NoError(err)

		expected := []models.Expense{
			{
				ExpenseID:   "gas bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(20),
				Name:        stringPtr("Gas Bill"),
				Notes:       "Monthly gas bill",
				CreatedDate: time.Date(2021, 9, 13, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "rent payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(800),
				Name:        stringPtr("Rent Payment"),
				Notes:       "Monthly rent payment",
				CreatedDate: time.Date(2021, 9, 12, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "student loan",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(150),
				Name:        stringPtr("Student Loan"),
				Notes:       "Monthly student loan payment",
				CreatedDate: time.Date(2021, 9, 11, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "credit card payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(200),
				Name:        stringPtr("Credit Card Payment"),
				Notes:       "Monthly credit card payment",
				CreatedDate: time.Date(2021, 9, 10, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "phone bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(25),
				Name:        stringPtr("Phone Bill"),
				Notes:       "Monthly phone bill",
				CreatedDate: time.Date(2021, 9, 9, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "water bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(30),
				Name:        stringPtr("Water Bill"),
				Notes:       "Monthly water bill",
				CreatedDate: time.Date(2021, 9, 8, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "internet bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(40),
				Name:        stringPtr("Internet Bill"),
				Notes:       "Monthly internet bill",
				CreatedDate: time.Date(2021, 9, 7, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "electricity bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(60),
				Name:        stringPtr("Electricity Bill"),
				Notes:       "Monthly electricity bill",
				CreatedDate: time.Date(2021, 9, 6, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "car insurance",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(75),
				Name:        stringPtr("Car Insurance"),
				Notes:       "Monthly car insurance premium",
				CreatedDate: time.Date(2021, 9, 5, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "spotify subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(9.99),
				Name:        stringPtr("Spotify Subscription"),
				Notes:       "Monthly Spotify subscription",
				CreatedDate: time.Date(2021, 9, 4, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "netflix subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(12.99),
				Name:        stringPtr("Netflix Subscription"),
				Notes:       "Monthly Netflix subscription",
				CreatedDate: time.Date(2021, 9, 3, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gym membership",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(50),
				Name:        stringPtr("Gym Membership"),
				Notes:       "Monthly gym membership",
				CreatedDate: time.Date(2021, 9, 2, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "insurance premium",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(100),
				Name:        stringPtr("Insurance Premium"),
				Notes:       "Monthly insurance premium",
				CreatedDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
		}

		for i, expense := range expensesResponse.Expenses {
			c.Equal(expected[i].ExpenseID, expense.ExpenseID)
		}
	})

	t.Run("Sorted by name - ASC", func(t *testing.T) {
		ger.ExpenseQueryParameters = &models.ExpenseQueryParameters{
			SortBy: "name",
			Period: periodID,
		}

		res, err := ger.GetByPeriod(ctx, apigwReq)
		c.NoError(err)
		c.Equal(res.StatusCode, http.StatusOK)

		var expensesResponse handlers.ExpensesResponse
		err = json.Unmarshal([]byte(res.Body), &expensesResponse)
		c.NoError(err)

		expected := []models.Expense{
			{
				ExpenseID:   "car insurance",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(75),
				Name:        stringPtr("Car Insurance"),
				Notes:       "Monthly car insurance premium",
				CreatedDate: time.Date(2021, 9, 5, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "credit card payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(200),
				Name:        stringPtr("Credit Card Payment"),
				Notes:       "Monthly credit card payment",
				CreatedDate: time.Date(2021, 9, 10, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "electricity bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(60),
				Name:        stringPtr("Electricity Bill"),
				Notes:       "Monthly electricity bill",
				CreatedDate: time.Date(2021, 9, 6, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gas bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(20),
				Name:        stringPtr("Gas Bill"),
				Notes:       "Monthly gas bill",
				CreatedDate: time.Date(2021, 9, 13, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gym membership",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(50),
				Name:        stringPtr("Gym Membership"),
				Notes:       "Monthly gym membership",
				CreatedDate: time.Date(2021, 9, 2, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "insurance premium",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(100),
				Name:        stringPtr("Insurance Premium"),
				Notes:       "Monthly insurance premium",
				CreatedDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "internet bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(40),
				Name:        stringPtr("Internet Bill"),
				Notes:       "Monthly internet bill",
				CreatedDate: time.Date(2021, 9, 7, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "netflix subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(12.99),
				Name:        stringPtr("Netflix Subscription"),
				Notes:       "Monthly Netflix subscription",
				CreatedDate: time.Date(2021, 9, 3, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "phone bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(25),
				Name:        stringPtr("Phone Bill"),
				Notes:       "Monthly phone bill",
				CreatedDate: time.Date(2021, 9, 9, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "rent payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(800),
				Name:        stringPtr("Rent Payment"),
				Notes:       "Monthly rent payment",
				CreatedDate: time.Date(2021, 9, 12, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "spotify subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(9.99),
				Name:        stringPtr("Spotify Subscription"),
				Notes:       "Monthly Spotify subscription",
				CreatedDate: time.Date(2021, 9, 4, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "student loan",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(150),
				Name:        stringPtr("Student Loan"),
				Notes:       "Monthly student loan payment",
				CreatedDate: time.Date(2021, 9, 11, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "water bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(30),
				Name:        stringPtr("Water Bill"),
				Notes:       "Monthly water bill",
				CreatedDate: time.Date(2021, 9, 8, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			}}

		for i, expense := range expensesResponse.Expenses {
			c.Equal(expected[i].ExpenseID, expense.ExpenseID)
		}
	})

	t.Run("Sorted by amount - ASC", func(t *testing.T) {
		ger.ExpenseQueryParameters = &models.ExpenseQueryParameters{
			SortBy: "amount",
			Period: periodID,
		}

		res, err := ger.GetByPeriod(ctx, apigwReq)
		c.NoError(err)
		c.Equal(res.StatusCode, http.StatusOK)

		var expensesResponse handlers.ExpensesResponse
		err = json.Unmarshal([]byte(res.Body), &expensesResponse)
		c.NoError(err)

		expected := []models.Expense{
			{
				ExpenseID:   "spotify subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(9.99),
				Name:        stringPtr("Spotify Subscription"),
				Notes:       "Monthly Spotify subscription",
				CreatedDate: time.Date(2021, 9, 4, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "netflix subscription",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(12.99),
				Name:        stringPtr("Netflix Subscription"),
				Notes:       "Monthly Netflix subscription",
				CreatedDate: time.Date(2021, 9, 3, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gas bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(20),
				Name:        stringPtr("Gas Bill"),
				Notes:       "Monthly gas bill",
				CreatedDate: time.Date(2021, 9, 13, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "phone bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(25),
				Name:        stringPtr("Phone Bill"),
				Notes:       "Monthly phone bill",
				CreatedDate: time.Date(2021, 9, 9, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "water bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(30),
				Name:        stringPtr("Water Bill"),
				Notes:       "Monthly water bill",
				CreatedDate: time.Date(2021, 9, 8, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "internet bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(40),
				Name:        stringPtr("Internet Bill"),
				Notes:       "Monthly internet bill",
				CreatedDate: time.Date(2021, 9, 7, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "gym membership",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(50),
				Name:        stringPtr("Gym Membership"),
				Notes:       "Monthly gym membership",
				CreatedDate: time.Date(2021, 9, 2, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "electricity bill",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(60),
				Name:        stringPtr("Electricity Bill"),
				Notes:       "Monthly electricity bill",
				CreatedDate: time.Date(2021, 9, 6, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "car insurance",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(75),
				Name:        stringPtr("Car Insurance"),
				Notes:       "Monthly car insurance premium",
				CreatedDate: time.Date(2021, 9, 5, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "insurance premium",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_1"),
				Amount:      floatPtr(100),
				Name:        stringPtr("Insurance Premium"),
				Notes:       "Monthly insurance premium",
				CreatedDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "student loan",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(150),
				Name:        stringPtr("Student Loan"),
				Notes:       "Monthly student loan payment",
				CreatedDate: time.Date(2021, 9, 11, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "credit card payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_2"),
				Amount:      floatPtr(200),
				Name:        stringPtr("Credit Card Payment"),
				Notes:       "Monthly credit card payment",
				CreatedDate: time.Date(2021, 9, 10, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
			{
				ExpenseID:   "rent payment",
				Username:    "e2e_test@gmail.com",
				CategoryID:  stringPtr("category_id_3"),
				Amount:      floatPtr(800),
				Name:        stringPtr("Rent Payment"),
				Notes:       "Monthly rent payment",
				CreatedDate: time.Date(2021, 9, 12, 0, 0, 0, 0, time.UTC),
				PeriodID:    "2021-09",
			},
		}

		for i, expense := range expensesResponse.Expenses {
			c.Equal(expected[i].ExpenseID, expense.ExpenseID)
		}
	})
}

func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
