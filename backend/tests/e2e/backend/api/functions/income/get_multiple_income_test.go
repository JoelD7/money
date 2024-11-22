package income

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/income/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	incomeTableName       string
	periodUserIncomeIndex string
	usersTableName        string
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	incomeTableName = env.GetString("INCOME_TABLE_NAME", "")
	periodUserIncomeIndex = env.GetString("PERIOD_USER_INCOME_INDEX", "")
	usersTableName = env.GetString("USERS_TABLE_NAME", "")

	os.Exit(m.Run())
}

func TestGetMultipleIncomeHandler(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	repo, err := income.NewDynamoRepository(dynamoClient, incomeTableName, periodUserIncomeIndex)
	c.Nil(err, "creating income repository failed")

	userRepo, err := users.NewDynamoRepository(dynamoClient, usersTableName)
	c.Nil(err, "creating users repository failed")

	cacheManager := cache.NewRedisCache()

	request := &handlers.GetMultipleIncomeRequest{
		IncomeRepo:   repo,
		CacheManager: cacheManager,
		Log:          logger.NewConsoleLogger("get_multiple_income_test"),
		Username:     setup.Username,
		PageSize:     10,
		StartKey:     "",
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": request.Username,
			},
		},
	}

	err = setup.CreateUser(ctx, userRepo, &models.User{}, t)
	c.Nil(err, "creating user failed")

	incomeEntries, err := setup.CreateIncomeEntries(ctx, repo, "", t)
	c.Nil(err, "creating income entries failed")
	c.NotEmptyf(incomeEntries, "income entries are empty")

	res, err := request.RouteToHandlers(ctx, apigwRequest)
	c.Nil(err, "get multiple income failed")

	var response handlers.MultipleIncomeResponse

	err = json.Unmarshal([]byte(res.Body), &response)
	c.Nil(err, "unmarshalling response failed")
	c.Len(response.Income, 10)
	c.NotEmpty(response.NextKey)
	c.Len(response.Periods, 6)
}

func BenchmarkGetMultipleIncome(t *testing.B) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	repo, err := income.NewDynamoRepository(dynamoClient, incomeTableName, periodUserIncomeIndex)
	c.Nil(err, "creating income repository failed")

	userRepo, err := users.NewDynamoRepository(dynamoClient, usersTableName)
	c.Nil(err, "creating users repository failed")

	cacheManager := cache.NewRedisCache()

	request := &handlers.GetMultipleIncomeRequest{
		IncomeRepo:   repo,
		CacheManager: cacheManager,
		Log:          logger.NewConsoleLogger("get_multiple_income_test"),
		Username:     setup.Username,
		PageSize:     10,
		StartKey:     "",
	}

	apigwRequest := &apigateway.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"username": request.Username,
			},
		},
	}

	err = setup.CreateUser(ctx, userRepo, &models.User{}, t)
	c.Nil(err, "creating user failed")

	incomeEntries, err := setup.CreateIncomeEntries(ctx, repo, "", t)
	c.Nil(err, "creating income entries failed")
	c.NotEmptyf(incomeEntries, "income entries are empty")

	for i := 0; i < t.N; i++ {
		res, err := request.RouteToHandlers(ctx, apigwRequest)
		c.Nil(err, "get multiple income failed")

		var response handlers.MultipleIncomeResponse

		err = json.Unmarshal([]byte(res.Body), &response)
		c.Nil(err, "unmarshalling response failed")
		c.Len(response.Income, 10)
		c.NotEmpty(response.NextKey)
		c.Len(response.Periods, 6)
	}
}
