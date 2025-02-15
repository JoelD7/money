package users

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/tests/e2e/api"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestCreateSavingGoals(t *testing.T) {
	c := require.New(t)

	requester, err := api.NewE2ERequester()
	c.Nil(err, "creating e2e requester failed")

	inputSavingGoal := new(models.SavingGoal)
	inputSavingGoal.SetName("test saving goal for e2e tests")
	inputSavingGoal.SetTarget(1000)
	inputSavingGoal.SetDeadline(time.Date(time.Now().Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC))

	createdSavingGoal, statusCode, err := requester.CreateSavingGoal(inputSavingGoal)
	c.Equal(http.StatusCreated, statusCode)
	c.Nil(err, "creating saving goal failed")
	c.NotNil(createdSavingGoal, "created saving goal is nil")
	c.NotEmpty(createdSavingGoal.GetSavingGoalID(), "created saving goal id is empty")

	t.Cleanup(func() {
		if createdSavingGoal.GetSavingGoalID() != "" {
			statusCode, err = requester.DeleteSavingGoal(createdSavingGoal.GetSavingGoalID())
			c.Equal(http.StatusNoContent, statusCode)
			c.Nil(err, "deleting saving goal failed")
		}
	})

	fetchedSavingGoal, statusCode, err := requester.GetSavingGoal(createdSavingGoal.GetSavingGoalID())
	c.Nil(err, "fetching saving goal failed")
	c.Equal(http.StatusOK, statusCode)
	c.NotNil(fetchedSavingGoal, "fetched saving goal is nil")
	c.Equal(createdSavingGoal.GetSavingGoalID(), fetchedSavingGoal.GetSavingGoalID(), "fetched saving goal id is different from the created one")
	c.Equal(inputSavingGoal.GetName(), fetchedSavingGoal.GetName(), "fetched saving goal name is different from the created one")
	c.Equal(inputSavingGoal.GetTarget(), fetchedSavingGoal.GetTarget(), "fetched saving goal target is different from the created one")
	c.Equal(inputSavingGoal.GetDeadline(), fetchedSavingGoal.GetDeadline(), "fetched saving goal deadline is different from the created one")
}

func TestSavingGoalsElimination(t *testing.T) {
	c := require.New(t)

	requester, err := api.NewE2ERequester()
	c.Nil(err, "creating e2e requester failed")

	inputSavingGoal := new(models.SavingGoal)
	inputSavingGoal.SetName("test saving goal for e2e tests")
	inputSavingGoal.SetTarget(1000)
	inputSavingGoal.SetDeadline(time.Date(time.Now().Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC))

	createdSavingGoal, statusCode, err := requester.CreateSavingGoal(inputSavingGoal)
	c.Equal(http.StatusCreated, statusCode)
	c.Nil(err, "creating saving goal failed")
	c.NotNil(createdSavingGoal, "created saving goal is nil")
	c.NotEmpty(createdSavingGoal.GetSavingGoalID(), "created saving goal id is empty")

	statusCode, err = requester.DeleteSavingGoal(createdSavingGoal.GetSavingGoalID())
	c.Equal(http.StatusNoContent, statusCode)
	c.Nil(err, "deleting saving goal failed")

	fetchedSavingGoal, statusCode, err := requester.GetSavingGoal(createdSavingGoal.GetSavingGoalID())
	c.Equal(http.StatusNotFound, statusCode)
	c.Nil(fetchedSavingGoal, "fetched saving goal is not nil after deletion")
	c.Contains(err.Error(), "Not found")
}
