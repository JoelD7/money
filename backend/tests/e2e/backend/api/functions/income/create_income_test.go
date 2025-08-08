package income

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestCreateIncome(t *testing.T) {
	c := require.New(t)

	requester, err := api.NewE2ERequester()
	c.Nil(err, "creating e2e requester failed")

	t.Run("Create failed: invalid period_id", func(t *testing.T) {
		period := &models.Period{
			Name:      aws.String("test period for e2e tests"),
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now().AddDate(0, 1, 0),
		}

		createdPeriod, statusCode, err := requester.CreatePeriod(period, t)
		c.Nil(err)
		c.Equal(http.StatusCreated, statusCode)
		c.NotNil(createdPeriod)

		income := &models.Income{
			Amount:   aws.Float64(1000),
			Name:     aws.String("test income for e2e tests"),
			PeriodID: stringPtr("random period id"),
		}

		createdIncome, statusCode, err := requester.CreateIncome(income, t)
		c.ErrorContains(err, "Invalid period")
		c.Equal(http.StatusBadRequest, statusCode)
		c.Nil(createdIncome)
	})
}

func stringPtr(s string) *string {
	return &s
}
