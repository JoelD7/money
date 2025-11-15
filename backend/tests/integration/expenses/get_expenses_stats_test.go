package expenses

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/api/functions/expenses/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetExpensesStats(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	username := "e2e_test@gmail.com"
	periodID := "2021-09"

	expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating expenses repository failed")

	usersRepo, err := users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
	c.Nil(err, "creating users repository failed")

	request := handlers.GetExpensesStatsRequest{
		ExpensesRepo: expensesRepo,
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
		CurrentPeriod: stringPtr(periodID),
	}

	err = setup.CreateUser(ctx, usersRepo, testUser, t)
	c.Nil(err, "creating user failed")

	_, err = setup.CreateExpensesEntries(ctx, expensesRepo, "", t)
	c.Nil(err, "creating expenses failed")

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
