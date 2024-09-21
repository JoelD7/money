package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/users/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestGetPeriodStats(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	username := "e2e_test@gmail.com"
	periodID := "2021-09"

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, expensesTableName, expensesRecurringTableName, periodUserExpenseIndex)
	c.Nil(err, "creating expenses repository failed")

	usersRepo, err := users.NewDynamoRepository(dynamoClient, usersTableName)
	c.Nil(err, "creating users repository failed")

	incomeRepo, err := income.NewDynamoRepository(dynamoClient, incomeTableName, periodUserIncomeIndex)
	c.Nil(err, "creating income repository failed")

	request := handlers.GetPeriodStatRequest{
		ExpensesRepo: expensesRepo,
		IncomeRepo:   incomeRepo,
		Log:          logger.NewConsoleLogger("get_period_stats_e2e_test"),
	}

	apigwRequest := &apigateway.Request{
		PathParameters: map[string]string{
			"periodID": periodID, //should be the same as in the sample json file
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": username,
			},
		},
	}

	testUser := &models.User{
		Username:      username,
		CurrentPeriod: periodID,
	}

	err = usersRepo.CreateUser(ctx, testUser)
	c.Nil(err, "creating user failed")

	defer t.Cleanup(func() {
		err = usersRepo.DeleteUser(ctx, username)
		c.Nil(err, "deleting user failed")
	})

	expensesList := setupExpenses(c)

	err = expensesRepo.BatchCreateExpenses(ctx, request.Log, expensesList)
	c.Nil(err, "batch creating expenses failed")

	defer t.Cleanup(func() {
		err = expensesRepo.BatchDeleteExpenses(ctx, expensesList)
		c.Nil(err, "batch deleting expenses failed")
	})

	incomeList := setupIncome(c)

	err = incomeRepo.BatchCreateIncome(ctx, incomeList)
	c.Nil(err, "batch creating income failed")

	defer t.Cleanup(func() {
		err = incomeRepo.BatchDeleteIncome(ctx, incomeList)
		c.Nil(err, "batch deleting income failed")
	})

	response, err := request.Process(ctx, apigwRequest)
	c.Nil(err, "get period stats failed")
	c.NotNil(response, "get period stats response is nil")
	c.Equal(http.StatusOK, response.StatusCode)

	var periodStat models.PeriodStat
	err = json.Unmarshal([]byte(response.Body), &periodStat)
	c.Nil(err, "unmarshalling response body failed")
	c.Len(periodStat.CategoryExpenseSummary, 3, "unexpected number of categories in the response")
	c.Equal(3000.00, periodStat.TotalIncome, fmt.Sprintf("expected %f, got %f", 3000.00, periodStat.TotalIncome))

	testValidatorByCategory := map[string]float64{
		"category_id_1": 172.98,
		"category_id_2": 430,
		"category_id_3": 970,
	}

	for _, summary := range periodStat.CategoryExpenseSummary {
		expected, ok := testValidatorByCategory[summary.CategoryID]
		c.True(ok, "unexpected category in the response")
		c.Equal(expected, summary.Total)
	}
}

func setupExpenses(c *require.Assertions) []*models.Expense {
	data, err := os.ReadFile("./samples/expenses_stats.json")
	c.Nil(err, "reading expenses sample file failed")

	var expensesList []*models.Expense
	err = json.Unmarshal(data, &expensesList)

	c.Len(expensesList, 13, "unexpected number of expenses in the sample file")

	return expensesList
}

func setupIncome(c *require.Assertions) []*models.Income {
	data, err := os.ReadFile("./samples/income_stats.json")
	c.Nil(err, "reading income sample file failed")

	var incomeList []*models.Income
	err = json.Unmarshal(data, &incomeList)

	c.Len(incomeList, 3, "unexpected number of income in the sample file")

	return incomeList
}
