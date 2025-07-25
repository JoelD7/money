package users

import (
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestDeletePeriod(t *testing.T) {
	c := require.New(t)

	t.Run("Delete period that doesn't exist", func(t *testing.T) {
		requester, err := api.NewE2ERequester()
		c.Nil(err, "creating e2e requester failed")

		statusCode, err := requester.DeletePeriod("random period id")
		c.Error(err)
		c.Equal(http.StatusNotFound, statusCode)
	})
}
