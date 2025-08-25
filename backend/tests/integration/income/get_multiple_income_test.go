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
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/JoelD7/money/backend/tests/e2e/setup"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	envConfig *models.EnvironmentConfiguration
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	envConfig = env.GetEnvConfig()

	logger.InitLogger(logger.ConsoleImplementation)

	os.Exit(m.Run())
}

func TestGetMultipleIncomeHandler(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	repo, err := income.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating income repository failed")

	userRepo, err := users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
	c.Nil(err, "creating users repository failed")

	periodRepo, err := period.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating period repository failed")

	cacheManager := cache.NewRedisCache()

	request := &handlers.GetMultipleIncomeRequest{
		IncomeRepo:   repo,
		CacheManager: cacheManager,
		PeriodRepo:   periodRepo,
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

	t.Run("Sort by name  ASC", func(t *testing.T) {
		request.QueryParameters = &models.QueryParameters{
			SortBy: "name",
			Period: "2023-7",
		}

		res, err = request.RouteToHandlers(ctx, apigwRequest)
		c.Nil(err, "get multiple income failed")

		err = json.Unmarshal([]byte(res.Body), &response)
		c.Nil(err, "unmarshalling response failed")

		expected := []*models.Income{
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN1lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(750),
				Name:       aws.String("income 1"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN2lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(600),
				Name:       aws.String("income 2"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN3lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(450),
				Name:       aws.String("income 3"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
		}

		for i, inc := range response.Income {
			c.Equal(expected[i].IncomeID, inc.IncomeID)
		}
	})

	t.Run("Sort by amount  ASC", func(t *testing.T) {
		request.QueryParameters = &models.QueryParameters{
			SortBy: "amount",
			Period: "2023-7",
		}

		res, err = request.RouteToHandlers(ctx, apigwRequest)
		c.Nil(err, "get multiple income failed")

		err = json.Unmarshal([]byte(res.Body), &response)
		c.Nil(err, "unmarshalling response failed")

		expected := []*models.Income{
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN3lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(450),
				Name:       aws.String("income 3"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN2lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(600),
				Name:       aws.String("income 2"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN1lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(750),
				Name:       aws.String("income 1"),
				PeriodID:   aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
		}

		for i, inc := range response.Income {
			c.Equal(expected[i].IncomeID, inc.IncomeID)
		}
	})
}

func BenchmarkGetMultipleIncome(t *testing.B) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	repo, err := income.NewDynamoRepository(dynamoClient, &models.EnvironmentConfiguration{IncomeTable: envConfig.IncomeTable, PeriodUserIncomeIndex: envConfig.PeriodUserIncomeIndex})
	c.Nil(err, "creating income repository failed")

	userRepo, err := users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
	c.Nil(err, "creating users repository failed")

	cacheManager := cache.NewRedisCache()

	request := &handlers.GetMultipleIncomeRequest{
		IncomeRepo:   repo,
		CacheManager: cacheManager,
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
