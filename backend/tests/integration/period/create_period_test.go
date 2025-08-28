package period

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/api/functions/users/handlers"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/storage/cache"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/JoelD7/money/backend/storage/savingoal"
	"github.com/JoelD7/money/backend/storage/savings"
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

// After adding new logic(mainly recurring saving generation) these tests are failing. I'll fix them later.
func TestProcess(t *testing.T) {
	c := require.New(t)
	sqsRetries := 3
	delay := time.Second * 1
	backoffFactor := 2

	t.Run("Set period to expenses without period", func(t *testing.T) {
		username := "e2e_test@gmail.com"

		apigwReq := &apigateway.Request{
			Headers: map[string]string{
				"Idempotency-Key": "1234",
			},
			Body: fmt.Sprintf(`{"username":"%s","name":"test-period","start_date":"2023-09-01T00:00:00Z","end_date":"2023-09-30T00:00:00Z"}`, username),
			RequestContext: events.APIGatewayProxyRequestContext{
				Authorizer: map[string]interface{}{
					"username": username,
				},
			},
		}

		ctx := context.Background()
		dynamoClient := dynamo.InitClient(ctx)

		periodRepo, err := period.NewDynamoRepository(dynamoClient, envConfig)
		c.Nil(err, "creating period repository failed")

		savingGoalRepo, err := savingoal.NewDynamoRepository(dynamoClient, envConfig)
		c.Nil(err)

		savingsRepo, err := savings.NewDynamoRepository(dynamoClient, envConfig)
		c.Nil(err)

		req := &handlers.CreatePeriodRequest{
			PeriodRepo:               periodRepo,
			IncomePeriodCacheManager: cache.NewRedisCache(),
			IdempotenceCache:         cache.NewRedisCache(),
			SavingGoalRepo:           savingGoalRepo,
			SavingsRepo:              savingsRepo,
		}

		expensesRepo, err := expenses.NewDynamoRepository(dynamoClient, envConfig)
		c.Nil(err, "creating expenses repository failed")

		expensesList, err := loadExpenses()
		c.Nil(err, "loading expenses from file failed")
		c.NotEmpty(username, "username from loaded expenses is empty")

		err = expensesRepo.BatchCreateExpenses(ctx, expensesList)
		c.Nil(err, "batch creating expenses failed")

		t.Cleanup(func() {
			err = expensesRepo.BatchDeleteExpenses(ctx, expensesList)
			c.Nil(err, "batch deleting expenses failed")

			p, err := req.PeriodRepo.GetLastPeriod(ctx, username)
			c.Nil(err, "couldn't delete created period: getting last period failed")

			err = req.PeriodRepo.DeletePeriod(ctx, p.ID, p.Username)
			c.Nil(err, "deleting period failed")
		})

		res, err := req.Process(ctx, apigwReq)
		c.Nil(err, "creating period failed")
		c.Equal(http.StatusCreated, res.StatusCode)

		var createdPeriod models.Period
		err = json.Unmarshal([]byte(res.Body), &createdPeriod)
		c.Nil(err, "unmarshalling created period failed")

		var expensesInPeriod []*models.Expense

		for i := 0; i < sqsRetries; i++ {
			//Wait for SQS to process the message
			time.Sleep(delay)

			expensesInPeriod, _, err = expensesRepo.GetExpensesByPeriod(ctx, createdPeriod.Username, &models.QueryParameters{Period: createdPeriod.ID, PageSize: 20})
			if expensesInPeriod != nil {
				break
			}

			delay *= time.Duration(backoffFactor)
		}

		c.Nil(err, "getting expenses by period failed")
		c.Len(expensesInPeriod, 18, fmt.Sprint("expected 18 expenses, got ", len(expensesInPeriod)))
	})
}

func TestCreatePeriod(t *testing.T) {
	c := require.New(t)

	// I do not recommend increasing these values. As of today I believe this a sufficiently large enough sample to test
	// what we need to test here. Increasing these values will make the tests slower, may produce an absurd amount of trash
	// logs and consume AWS free-tier resources best used for the real application.
	numSavingGoals := 12
	numPeriods := 3

	t.Run("Create recurring savings across multiple periods", func(t *testing.T) {
		requester, err := api.NewE2ERequester()
		c.NoError(err)

		savingGoals := make([]models.SavingGoal, 0, numSavingGoals)
		periods := make([]models.Period, 0, numPeriods)
		savingsByPeriod := make(map[string][]*models.Saving)

		// Step 1: Create 12 recurring saving goals
		for i := 0; i < numSavingGoals; i++ {
			name := "Saving Goal " + string(rune('A'+i))
			target := float64(1000 * (i + 1))
			recurringAmount := float64(100 * (i + 1))

			savingGoal := &models.SavingGoal{
				Name:            &name,
				Target:          &target,
				Progress:        new(float64),
				IsRecurring:     true,
				RecurringAmount: &recurringAmount,
			}

			createdSavingGoal, statusCode, err := requester.CreateSavingGoal(savingGoal, t)
			c.NoError(err)
			c.Equal(http.StatusCreated, statusCode)
			c.NotEmpty(createdSavingGoal.SavingGoalID)

			savingGoals = append(savingGoals, *createdSavingGoal)
		}

		var statusCode int
		var createdPeriod *models.Period
		var periodToCreate *models.Period
		var startDate, endDate time.Time
		var periodName string

		// Step 2: Create 3 periods
		currentDate := time.Now()
		for i := 0; i < numPeriods; i++ {
			periodName = "Period " + string(rune('A'+i))

			// Each period is one month
			startDate = currentDate.AddDate(0, i, 0)
			endDate = currentDate.AddDate(0, i+1, -1)

			periodToCreate = &models.Period{
				Name:      &periodName,
				StartDate: startDate,
				EndDate:   endDate,
			}

			createdPeriod, statusCode, err = requester.CreatePeriod(periodToCreate, t)
			derefCreatedPeriod := *createdPeriod //outside of the deferred call
			// Creating periods automatically creates savings, so this cleanup is necessary
			t.Cleanup(func() {
				periodSavings := savingsByPeriod[*derefCreatedPeriod.Name]
				if len(periodSavings) > 0 {
					for _, saving := range periodSavings {
						statusCode, err = requester.DeleteSaving(saving.SavingID)
						c.NoError(err)
						c.Equal(http.StatusNoContent, statusCode)
						delete(savingsByPeriod, *derefCreatedPeriod.Name)
					}
				}
			})

			c.NoError(err)
			c.Equal(http.StatusCreated, statusCode)
			c.NotEmpty(createdPeriod.ID)

			periods = append(periods, *createdPeriod)
		}

		var params *models.QueryParameters
		var savings []*models.Saving

		// Step 3: Verify that savings were created for each recurring saving goal in each period
		for _, period := range periods {

			c.NotNil(period.Name)
			params = &models.QueryParameters{
				Period:   *period.Name,
				PageSize: 100, // Set a large enough page size to get all savings
			}

			// Get savings for this period
			savings, _, statusCode, err = requester.GetSavings(params)
			c.NoError(err)
			savingsByPeriod[*period.Name] = append(savingsByPeriod[*period.Name], savings...)
			c.Equal(200, statusCode)

			// We should have one saving per saving goal for each period
			c.Len(savings, numSavingGoals)

			// Create a map to track which saving goals have been found
			foundSavingGoals := make(map[string]bool)
			for _, savingGoal := range savingGoals {
				foundSavingGoals[savingGoal.SavingGoalID] = false
			}

			// Verify each saving matches a saving goal and has the correct amount
			for _, saving := range savings {
				// Find the corresponding saving goal
				var matchingSavingGoal models.SavingGoal
				found := false
				for _, savingGoal := range savingGoals {
					if *saving.SavingGoalID == savingGoal.SavingGoalID {
						matchingSavingGoal = savingGoal
						found = true
						break
					}
				}

				c.True(found, "Saving corresponds to a saving goal")

				// Mark this saving goal as found
				foundSavingGoals[matchingSavingGoal.SavingGoalID] = true

				// Verify saving has correct period
				c.NotNil(saving.PeriodID)
				c.NotNil(period.Name)
				c.Equal(*period.Name, *saving.PeriodID)

				// Verify amount matches the recurring amount from the saving goal
				c.Equal(*matchingSavingGoal.RecurringAmount, *saving.Amount)
				c.NotNil(matchingSavingGoal.Name)
				c.Equal(*matchingSavingGoal.Name, saving.SavingGoalName)
			}

			// Verify all saving goals were found
			for savingGoalID, found := range foundSavingGoals {
				c.True(found, "Saving goal "+savingGoalID+" should have a saving in period "+period.ID)
			}
		}
	})

	t.Run("Non-recurring saving goals don't create automatic savings", func(t *testing.T) {
		requester, err := api.NewE2ERequester()
		c.NoError(err)

		// Create a non-recurring saving goal
		goalName := "One-time Purchase"
		target := float64(500)

		savingGoal := &models.SavingGoal{
			Name:        &goalName,
			Target:      &target,
			Progress:    new(float64), // Initially 0
			IsRecurring: false,
		}

		createdSavingGoal, statusCode, err := requester.CreateSavingGoal(savingGoal, t)
		c.NoError(err)
		c.Equal(http.StatusCreated, statusCode)
		c.NotEmpty(createdSavingGoal.SavingGoalID)

		// Create a period
		periodName := "February 2025"
		startDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 2, 28, 23, 59, 59, 0, time.UTC)

		period := models.Period{
			Name:      &periodName,
			StartDate: startDate,
			EndDate:   endDate,
		}

		createdPeriod, statusCode, err := requester.CreatePeriod(&period, t)
		c.NoError(err)
		c.Equal(http.StatusCreated, statusCode)
		c.NotEmpty(createdPeriod.ID)

		// Get savings for this period and saving goal
		params := &models.QueryParameters{
			Period:       *createdPeriod.Name,
			SavingGoalID: createdSavingGoal.SavingGoalID,
		}

		savings, _, statusCode, err := requester.GetSavings(params)
		c.NotNil(err)
		c.Equal(http.StatusNotFound, statusCode)

		// Verify no savings were automatically created
		c.Len(savings, 0)
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
