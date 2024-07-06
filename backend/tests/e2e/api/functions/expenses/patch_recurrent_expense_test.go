package expenses

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	expensesRepo "github.com/JoelD7/money/backend/storage/expenses"
	periodRepo "github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/tests/e2e/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	c := require.New(t)

	dynamoClient := utils.InitDynamoClient()
	ctx := context.Background()

	req := &handlers.PatchRecurrentExpenseRequest{
		ExpensesRepo: expensesRepo.NewDynamoRepository(dynamoClient),
		PeriodRepo:   periodRepo.NewDynamoRepository(dynamoClient),
		Log:          logger.NewConsoleLogger("patch_recurrent_expense_e2e_test"),
	}

	username := "e2e_test@gmail.com"

	expensesList, err := loadExpenses()
	c.Nil(err, "loading expenses from file failed")
	c.NotEmpty(username, "username from loaded expenses is empty")

	startDate, err := time.Parse(time.DateOnly, "2023-09-01")
	c.Nil(err, "parsing start date failed")

	endDate, err := time.Parse(time.DateOnly, "2023-09-30")
	c.Nil(err, "parsing end date failed")

	period := &models.Period{
		ID:        "test-period",
		Name:      utils.StringPtr("test-period"),
		Username:  username,
		StartDate: startDate,
		EndDate:   endDate,
	}

	_, err = req.PeriodRepo.CreatePeriod(ctx, period)
	c.Nil(err, "creating period failed")

	apigwReq := &apigateway.Request{
		Body: fmt.Sprintf(`{"period": "%s"}`, period.ID),
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": username,
			},
		},
	}

	err = req.ExpensesRepo.BatchCreateExpenses(ctx, req.Log, expensesList)
	c.Nil(err, "batch creating expenses failed")

	t.Cleanup(func() {
		err = req.ExpensesRepo.BatchDeleteExpenses(ctx, expensesList)
		c.Nil(err, "batch deleting expenses failed")

		err = req.PeriodRepo.DeletePeriod(ctx, period.ID, period.Username)
		c.Nil(err, "deleting period failed")
	})

	res, err := req.Process(ctx, apigwReq)
	c.Nil(err)
	c.Equal(http.StatusOK, res.StatusCode)

	result, _, err := req.ExpensesRepo.GetExpensesByPeriod(ctx, period.Username, period.ID, "", 20)
	c.Nil(err, "getting expenses by period failed")
	c.Len(result, 18, fmt.Sprint("expected 18 expenses, got ", len(result)))
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
