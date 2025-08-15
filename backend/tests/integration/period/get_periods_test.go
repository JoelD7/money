package period

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/period"
	"github.com/stretchr/testify/require"
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

func TestGetPeriods(t *testing.T) {
	c := require.New(t)

	now := time.Now()
	Username := "e2e_test@gmail.com"
	ctx := context.Background()

	dynamoClient := dynamo.InitClient(ctx)

	periodRepo, err := period.NewDynamoRepository(dynamoClient, envConfig)
	c.Nil(err, "creating period repository failed")

	//Create periods, 13 past, 7 within and 17 future
	periodsToCreate := []*models.Period{
		// 13 Past Periods (StartDate and EndDate are before now)
		{
			Username:  Username,
			Name:      stringPtr("Past Period 1"),
			StartDate: now.AddDate(0, -13, 0),
			EndDate:   now.AddDate(0, -13, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 2"),
			StartDate: now.AddDate(0, -12, 0),
			EndDate:   now.AddDate(0, -12, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 3"),
			StartDate: now.AddDate(0, -11, 0),
			EndDate:   now.AddDate(0, -11, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 4"),
			StartDate: now.AddDate(0, -10, 0),
			EndDate:   now.AddDate(0, -10, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 5"),
			StartDate: now.AddDate(0, -9, 0),
			EndDate:   now.AddDate(0, -9, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 6"),
			StartDate: now.AddDate(0, -8, 0),
			EndDate:   now.AddDate(0, -8, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 7"),
			StartDate: now.AddDate(0, -7, 0),
			EndDate:   now.AddDate(0, -7, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 8"),
			StartDate: now.AddDate(0, -6, 0),
			EndDate:   now.AddDate(0, -6, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 9"),
			StartDate: now.AddDate(0, -5, 0),
			EndDate:   now.AddDate(0, -5, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 10"),
			StartDate: now.AddDate(0, -4, 0),
			EndDate:   now.AddDate(0, -4, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 11"),
			StartDate: now.AddDate(0, -3, 0),
			EndDate:   now.AddDate(0, -3, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 12"),
			StartDate: now.AddDate(0, -2, 0),
			EndDate:   now.AddDate(0, -2, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Past Period 13"),
			StartDate: now.AddDate(0, -1, 0),
			EndDate:   now.AddDate(0, -1, 28),
		},

		// 7 Active Periods (StartDate <= now < EndDate)
		{
			Username:  Username,
			Name:      stringPtr("Active Period 1"),
			StartDate: now.AddDate(0, 0, -15),
			EndDate:   now.AddDate(0, 0, 15),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 2"),
			StartDate: now.AddDate(0, 0, -10),
			EndDate:   now.AddDate(0, 0, 20),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 3"),
			StartDate: now.AddDate(0, -1, 0),
			EndDate:   now.AddDate(0, 1, 0),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 4 (Starts now)"),
			StartDate: now,
			EndDate:   now.AddDate(0, 1, 0),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 5"),
			StartDate: now.AddDate(0, -2, 0),
			EndDate:   now.AddDate(0, 2, 0),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 6"),
			StartDate: now.AddDate(0, 0, -5),
			EndDate:   now.AddDate(0, 0, 5),
		},
		{
			Username:  Username,
			Name:      stringPtr("Active Period 7"),
			StartDate: now.AddDate(0, 0, -1),
			EndDate:   now.AddDate(0, 0, 1),
		},

		// 17 Future Periods (StartDate and EndDate are after now)
		{
			Username:  Username,
			Name:      stringPtr("Future Period 1"),
			StartDate: now.AddDate(0, 1, 0),
			EndDate:   now.AddDate(0, 1, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 2"),
			StartDate: now.AddDate(0, 2, 0),
			EndDate:   now.AddDate(0, 2, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 3"),
			StartDate: now.AddDate(0, 3, 0),
			EndDate:   now.AddDate(0, 3, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 4"),
			StartDate: now.AddDate(0, 4, 0),
			EndDate:   now.AddDate(0, 4, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 5"),
			StartDate: now.AddDate(0, 5, 0),
			EndDate:   now.AddDate(0, 5, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 6"),
			StartDate: now.AddDate(0, 6, 0),
			EndDate:   now.AddDate(0, 6, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 7"),
			StartDate: now.AddDate(0, 7, 0),
			EndDate:   now.AddDate(0, 7, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 8"),
			StartDate: now.AddDate(0, 8, 0),
			EndDate:   now.AddDate(0, 8, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 9"),
			StartDate: now.AddDate(0, 9, 0),
			EndDate:   now.AddDate(0, 9, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 10"),
			StartDate: now.AddDate(0, 10, 0),
			EndDate:   now.AddDate(0, 10, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 11"),
			StartDate: now.AddDate(0, 11, 0),
			EndDate:   now.AddDate(0, 11, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 12"),
			StartDate: now.AddDate(0, 12, 0),
			EndDate:   now.AddDate(0, 12, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 13"),
			StartDate: now.AddDate(0, 13, 0),
			EndDate:   now.AddDate(0, 13, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 14"),
			StartDate: now.AddDate(0, 14, 0),
			EndDate:   now.AddDate(0, 14, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 15"),
			StartDate: now.AddDate(0, 15, 0),
			EndDate:   now.AddDate(0, 15, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 16"),
			StartDate: now.AddDate(0, 16, 0),
			EndDate:   now.AddDate(0, 16, 28),
		},
		{
			Username:  Username,
			Name:      stringPtr("Future Period 17"),
			StartDate: now.AddDate(0, 17, 0),
			EndDate:   now.AddDate(0, 17, 28),
		},
	}

	for _, period := range periodsToCreate {
		createdPeriod, err := periodRepo.CreatePeriod(ctx, period)
		c.Nil(err)
		c.NotNil(createdPeriod)
	}

	active := true
	pageSize := 10

	periods, nextKey, err := periodRepo.GetPeriods(ctx, Username, "", pageSize, active)
	c.Nil(err)
	c.NotEmpty(nextKey)
	c.Len(periods, pageSize)
	c.True(arePeriodsSorted(periods, true))
	c.True(arePeriodsActive(periods))

	periods, nextKey, err = periodRepo.GetPeriods(ctx, Username, nextKey, pageSize, active)
	c.Nil(err)
	c.NotEmpty(nextKey)
	c.Len(periods, pageSize)
	c.True(arePeriodsSorted(periods, true))
	c.True(arePeriodsActive(periods))

	periods, nextKey, err = periodRepo.GetPeriods(ctx, Username, nextKey, pageSize, active)
	c.Nil(err)
	c.Empty(nextKey)
	c.Len(periods, 4)
	c.True(arePeriodsSorted(periods, true))
	c.True(arePeriodsActive(periods))
}

func arePeriodsSorted(periods []*models.Period, asc bool) bool {
	for i := 0; i < len(periods)-1; i++ {
		if asc {
			if periods[i].EndDate.After(periods[i+1].EndDate) {
				return false
			}
		} else {
			if periods[i].EndDate.Before(periods[i+1].EndDate) {
				return false
			}
		}
	}

	return true
}

func arePeriodsActive(periods []*models.Period) bool {
	now := time.Now()

	for _, p := range periods {
		if !(isBetweenDates(now, p.StartDate, p.EndDate) || areDatesAfter(now, p.StartDate, p.EndDate)) {
			return false
		}
	}

	return true
}

func areDatesAfter(date time.Time, startDate time.Time, endDate time.Time) bool {
	return startDate.After(date) && endDate.After(date)
}

func isBetweenDates(date time.Time, startDate time.Time, endDate time.Time) bool {
	return (startDate.Before(date) || startDate.Equal(date)) && (endDate.After(date) || endDate.Equal(date))
}

func stringPtr(s string) *string {
	return &s
}
