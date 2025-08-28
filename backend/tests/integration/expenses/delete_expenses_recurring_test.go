package expenses

import (
	"context"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	expensesRecurring "github.com/JoelD7/money/backend/storage/expenses-recurring"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

var (
	expensesRecurringTableName string
)

func TestProcess(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	username := "e2e_test@gmail.com"

	ex := &models.Expense{
		ExpenseID:    "test-expense-id",
		Username:     username,
		Amount:       aws.Float64(150.34),
		RecurringDay: aws.Int(10),
		IsRecurring:  true,
		Name:         aws.String("Test Expense"),
	}

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating expenses repository failed")

	createdExpense, err := expensesRepo.CreateExpense(ctx, ex)
	c.Nil(err, "creating expense failed")
	c.NotNil(createdExpense, "created expense is nil")

	t.Cleanup(func() {
		err = expensesRepo.DeleteExpense(ctx, createdExpense.ExpenseID, username)
		c.Nil(err, "deleting expense failed")
	})

	expensesRecRepo, err := expensesRecurring.NewExpenseRecurringDynamoRepository(dynamoClient, expensesRecurringTableName)
	c.Nil(err, "creating expenses recurring repository failed")

	req := &handlers.DeleteExpenseRecurringRequest{
		Repo: expensesRecRepo,
	}

	expenseRecurringID := strings.ToLower(*createdExpense.Name)

	apigwReq := &apigateway.Request{
		PathParameters: map[string]string{
			"expenseRecurringID": expenseRecurringID,
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": username,
			},
		},
	}

	res, err := req.Process(ctx, apigwReq)
	c.Nil(err, "deleting expense recurring failed")
	c.Equal(http.StatusNoContent, res.StatusCode)

	_, err = expensesRepo.GetExpense(ctx, username, createdExpense.ExpenseID)
	c.Nil(err, "getting expense failed")

	_, err = req.Repo.GetExpenseRecurring(ctx, expenseRecurringID, username)
	c.ErrorIs(err, models.ErrRecurringExpenseNotFound)
}
