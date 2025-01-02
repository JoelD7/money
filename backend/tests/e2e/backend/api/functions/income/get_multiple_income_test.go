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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	incomeTableName             string
	periodUserIncomeIndex       string
	usersTableName              string
	periodUserNameIncomeIDIndex string
	periodUserAmountIndex       string
	periodUserCreatedDateIndex  string
	usernameCreatedDateIndex    string
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	incomeTableName = env.GetString("INCOME_TABLE_NAME", "")
	periodUserIncomeIndex = env.GetString("PERIOD_USER_INCOME_INDEX", "")
	usersTableName = env.GetString("USERS_TABLE_NAME", "")
	periodUserNameIncomeIDIndex = env.GetString("PERIOD_USER_NAME_INCOME_ID_INDEX", "")
	periodUserAmountIndex = env.GetString("PERIOD_USER_AMOUNT_INDEX", "")
	periodUserCreatedDateIndex = env.GetString("PERIOD_USER_CREATED_DATE_INDEX", "")
	usernameCreatedDateIndex = env.GetString("USERNAME_CREATED_DATE_INDEX", "")

	os.Exit(m.Run())
}

func TestGetMultipleIncomeHandler(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	dynamoClient := dynamo.InitClient(ctx)

	envConfig := &models.EnvironmentConfiguration{
		IncomeTable:                 incomeTableName,
		PeriodUserIncomeIndex:       periodUserIncomeIndex,
		PeriodUserNameIncomeIDIndex: periodUserNameIncomeIDIndex,
		PeriodUserAmountIndex:       periodUserAmountIndex,
		PeriodUserCreatedDateIndex:  periodUserCreatedDateIndex,
		UsernameCreatedDateIndex:    usernameCreatedDateIndex,
	}

	repo, err := income.NewDynamoRepository(dynamoClient, envConfig)
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
				Period:     aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN2lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(600),
				Name:       aws.String("income 2"),
				Period:     aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN3lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(450),
				Name:       aws.String("income 3"),
				Period:     aws.String("2023-7"),
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
				Period:     aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN2lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(600),
				Name:       aws.String("income 2"),
				Period:     aws.String("2023-7"),
				PeriodUser: aws.String("2023-7:test@gmail.com"),
			},
			{
				Username:   "e2e_test@gmail.com",
				IncomeID:   "IN1lVnB1tCaSpyQLDEUswM",
				Amount:     aws.Float64(750),
				Name:       aws.String("income 1"),
				Period:     aws.String("2023-7"),
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

	repo, err := income.NewDynamoRepository(dynamoClient, &models.EnvironmentConfiguration{IncomeTable: incomeTableName, PeriodUserIncomeIndex: periodUserIncomeIndex})
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
