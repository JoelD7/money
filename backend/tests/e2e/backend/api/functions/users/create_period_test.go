package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/users/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func init() {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(err)
	}
}

func TestProcess(t *testing.T) {
	c := require.New(t)

	t.Run("Set period to expenses without period", func(t *testing.T) {
		username := "e2e_test@gmail.com"

		apigwReq := &apigateway.Request{
			Body: fmt.Sprintf(`{"username":"%s","name":"test-period","start_date":"2023-09-01T00:00:00Z","end_date":"2023-09-30T00:00:00Z"}`, username),
			RequestContext: events.APIGatewayProxyRequestContext{
				Authorizer: map[string]interface{}{
					"username": username,
				},
			},
		}

		ctx := context.Background()
		dynamoClient := dynamo.InitDynamoClient(ctx)

		req := &handlers.CreatePeriodRequest{
			PeriodRepo: period.NewDynamoRepository(dynamoClient),
			Log:        logger.NewConsoleLogger("create_period_e2e_test"),
		}

		expensesRepo := expenses.NewDynamoRepository(dynamoClient)

		expensesList, err := loadExpenses()
		c.Nil(err, "loading expenses from file failed")
		c.NotEmpty(username, "username from loaded expenses is empty")

		err = expensesRepo.BatchCreateExpenses(ctx, req.Log, expensesList)
		c.Nil(err, "batch creating expenses failed")

		t.Cleanup(func() {
			err = expensesRepo.BatchDeleteExpenses(ctx, expensesList)
			c.Nil(err, "batch deleting expenses failed")

			p, err := req.PeriodRepo.GetLastPeriod(ctx, username)
			c.Nil(err, "getting last period failed")

			err = req.PeriodRepo.DeletePeriod(ctx, p.ID, p.Username)
			c.Nil(err, "deleting period failed")
		})

		res, err := req.Process(ctx, apigwReq)
		c.Nil(err, "creating period failed")
		c.Equal(http.StatusCreated, res.StatusCode)

		var createdPeriod models.Period
		err = json.Unmarshal([]byte(res.Body), &createdPeriod)
		c.Nil(err, "unmarshalling created period failed")

		result, _, err := expensesRepo.GetExpensesByPeriod(ctx, createdPeriod.Username, createdPeriod.ID, "", 20)
		c.Nil(err, "getting expenses by period failed")
		c.Len(result, 18, fmt.Sprint("expected 18 expenses, got ", len(result)))
	})
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
