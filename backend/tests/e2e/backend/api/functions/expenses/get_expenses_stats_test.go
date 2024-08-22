package expenses

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestGetExpensesStats(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	username := "e2e_test@gmail.com"

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, expensesTableName, expensesRecurringTableName, periodUserIndex)
	c.Nil(err, "creating expenses repository failed")

	request := handlers.GetExpensesStatsRequest{
		ExpensesRepo: expensesRepo,
		Log:          logger.NewConsoleLogger("get_expenses_stats_e2e_test"),
	}

	apigwRequest := &apigateway.Request{
		PathParameters: map[string]string{
			"periodID": "2021-09", //should be the same as in the sample json file
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": username,
			},
		},
	}

	expensesList := setupExpenses(c)

	err = expensesRepo.BatchCreateExpenses(ctx, request.Log, expensesList)
	c.Nil(err, "batch creating expenses failed")

	t.Cleanup(func() {
		err = expensesRepo.BatchDeleteExpenses(ctx, expensesList)
		c.Nil(err, "batch deleting expenses failed")
	})

	response, err := request.Process(ctx, apigwRequest)
	c.Nil(err, "get expenses stats failed")
	c.NotNil(response, "get expenses stats response is nil")
	c.Equal(http.StatusOK, response.StatusCode)

	var categoryExpenseSummary []*models.CategoryExpenseSummary
	err = json.Unmarshal([]byte(response.Body), &categoryExpenseSummary)
	c.Nil(err, "unmarshalling response body failed")
	c.Len(categoryExpenseSummary, 3, "unexpected number of categories in the response")

	testValidatorByCategory := map[string]float64{
		"category_id_1": 172.98,
		"category_id_2": 430,
		"category_id_3": 970,
	}

	for _, summary := range categoryExpenseSummary {
		expected, ok := testValidatorByCategory[summary.CategoryID]
		c.True(ok, "unexpected category in the response")
		c.Equal(expected, summary.Total)
	}
}

func setupExpenses(c *require.Assertions) []*models.Expense {
	data, err := os.ReadFile("./samples/expenses.json")
	c.Nil(err, "reading expenses sample file failed")

	var expensesList []*models.Expense
	err = json.Unmarshal(data, &expensesList)

	c.Len(expensesList, 13, "unexpected number of expenses in the sample file")

	return expensesList
}
