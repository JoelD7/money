package users

import (
	"net/http"
	"testing"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/stretchr/testify/require"
)

func TestGetPeriods(t *testing.T) {
	c := require.New(t)

	requester, err := api.NewE2ERequester()
	c.Nil(err)

	//setup
	now := time.Now()

	//Create periods, 13 past, 7 within and 17 future
	periodsToCreate := []*models.Period{
		// 13 Past Periods (StartDate and EndDate are before now)
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 1"),
			StartDate: now.AddDate(0, -13, 0),
			EndDate:   now.AddDate(0, -13, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 2"),
			StartDate: now.AddDate(0, -12, 0),
			EndDate:   now.AddDate(0, -12, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 3"),
			StartDate: now.AddDate(0, -11, 0),
			EndDate:   now.AddDate(0, -11, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 4"),
			StartDate: now.AddDate(0, -10, 0),
			EndDate:   now.AddDate(0, -10, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 5"),
			StartDate: now.AddDate(0, -9, 0),
			EndDate:   now.AddDate(0, -9, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 6"),
			StartDate: now.AddDate(0, -8, 0),
			EndDate:   now.AddDate(0, -8, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 7"),
			StartDate: now.AddDate(0, -7, 0),
			EndDate:   now.AddDate(0, -7, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 8"),
			StartDate: now.AddDate(0, -6, 0),
			EndDate:   now.AddDate(0, -6, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 9"),
			StartDate: now.AddDate(0, -5, 0),
			EndDate:   now.AddDate(0, -5, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 10"),
			StartDate: now.AddDate(0, -4, 0),
			EndDate:   now.AddDate(0, -4, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 11"),
			StartDate: now.AddDate(0, -3, 0),
			EndDate:   now.AddDate(0, -3, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 12"),
			StartDate: now.AddDate(0, -2, 0),
			EndDate:   now.AddDate(0, -2, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Past Period 13"),
			StartDate: now.AddDate(0, -1, 0),
			EndDate:   now.AddDate(0, -1, 28),
		},

		// 7 Active Periods (StartDate <= now < EndDate)
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 1"),
			StartDate: now.AddDate(0, 0, -15),
			EndDate:   now.AddDate(0, 0, 15),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 2"),
			StartDate: now.AddDate(0, 0, -10),
			EndDate:   now.AddDate(0, 0, 20),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 3"),
			StartDate: now.AddDate(0, -1, 0),
			EndDate:   now.AddDate(0, 1, 0),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 4 (Starts now)"),
			StartDate: now,
			EndDate:   now.AddDate(0, 1, 0),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 5"),
			StartDate: now.AddDate(0, -2, 0),
			EndDate:   now.AddDate(0, 2, 0),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 6"),
			StartDate: now.AddDate(0, 0, -5),
			EndDate:   now.AddDate(0, 0, 5),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Active Period 7"),
			StartDate: now.AddDate(0, 0, -1),
			EndDate:   now.AddDate(0, 0, 1),
		},

		// 17 Future Periods (StartDate and EndDate are after now)
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 1"),
			StartDate: now.AddDate(0, 1, 0),
			EndDate:   now.AddDate(0, 1, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 2"),
			StartDate: now.AddDate(0, 2, 0),
			EndDate:   now.AddDate(0, 2, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 3"),
			StartDate: now.AddDate(0, 3, 0),
			EndDate:   now.AddDate(0, 3, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 4"),
			StartDate: now.AddDate(0, 4, 0),
			EndDate:   now.AddDate(0, 4, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 5"),
			StartDate: now.AddDate(0, 5, 0),
			EndDate:   now.AddDate(0, 5, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 6"),
			StartDate: now.AddDate(0, 6, 0),
			EndDate:   now.AddDate(0, 6, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 7"),
			StartDate: now.AddDate(0, 7, 0),
			EndDate:   now.AddDate(0, 7, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 8"),
			StartDate: now.AddDate(0, 8, 0),
			EndDate:   now.AddDate(0, 8, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 9"),
			StartDate: now.AddDate(0, 9, 0),
			EndDate:   now.AddDate(0, 9, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 10"),
			StartDate: now.AddDate(0, 10, 0),
			EndDate:   now.AddDate(0, 10, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 11"),
			StartDate: now.AddDate(0, 11, 0),
			EndDate:   now.AddDate(0, 11, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 12"),
			StartDate: now.AddDate(0, 12, 0),
			EndDate:   now.AddDate(0, 12, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 13"),
			StartDate: now.AddDate(0, 13, 0),
			EndDate:   now.AddDate(0, 13, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 14"),
			StartDate: now.AddDate(0, 14, 0),
			EndDate:   now.AddDate(0, 14, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 15"),
			StartDate: now.AddDate(0, 15, 0),
			EndDate:   now.AddDate(0, 15, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 16"),
			StartDate: now.AddDate(0, 16, 0),
			EndDate:   now.AddDate(0, 16, 28),
		},
		{
			Username:  requester.Username,
			Name:      stringPtr("Future Period 17"),
			StartDate: now.AddDate(0, 17, 0),
			EndDate:   now.AddDate(0, 17, 28),
		},
	}

	for _, period := range periodsToCreate {
		createdPeriod, statusCode, err := requester.CreatePeriod(period, t)
		c.Nil(err)
		c.Equal(http.StatusCreated, statusCode)
		c.NotNil(createdPeriod)
	}

	t.Run("Get active periods", func(t *testing.T) {
		active := true
		pageSize := 10

		res, statusCode, err := requester.GetPeriods(active, &models.QueryParameters{PageSize: pageSize})
		c.Nil(err)
		c.Equal(http.StatusOK, statusCode)
		c.NotNil(res)
		c.NotEmpty(res.NextKey)
		c.Len(res.Periods, pageSize)
		c.True(arePeriodsSortedByEndDate(res.Periods, true))
		c.True(arePeriodsActive(res.Periods))

		res, statusCode, err = requester.GetPeriods(active, &models.QueryParameters{PageSize: pageSize, StartKey: res.NextKey})
		c.Nil(err)
		c.Equal(http.StatusOK, statusCode)
		c.NotNil(res)
		c.NotEmpty(res.NextKey)
		c.Len(res.Periods, pageSize)
		c.True(arePeriodsSortedByEndDate(res.Periods, true))
		c.True(arePeriodsActive(res.Periods))

		res, statusCode, err = requester.GetPeriods(active, &models.QueryParameters{PageSize: pageSize, StartKey: res.NextKey})
		c.Nil(err)
		c.Equal(http.StatusOK, statusCode)
		c.NotNil(res)
		c.Empty(res.NextKey) //This is the last page
		c.Len(res.Periods, 4)
		c.True(arePeriodsSortedByEndDate(res.Periods, true))
		c.True(arePeriodsActive(res.Periods))
	})

	t.Run("Get all periods", func(t *testing.T) {
		active := false
		pageSize := 40

		res, statusCode, err := requester.GetPeriods(active, &models.QueryParameters{PageSize: pageSize})
		c.Nil(err)
		c.Equal(http.StatusOK, statusCode)
		c.NotNil(res)
		c.Empty(res.NextKey)
		c.Len(res.Periods, 37)
	})
}

func arePeriodsSortedByEndDate(periods []*models.Period, asc bool) bool {
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
