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

func TestGetSavingGoals(t *testing.T) {
	c := require.New(t)

	requester, err := api.NewE2ERequester()
	c.NoError(err, "creating e2e requester failed")

	var createdGoals []*models.SavingGoal
	var goalIDs []string

	createGoal := func(name string, target float64, daysFromNow int) *models.SavingGoal {
		inputGoal := new(models.SavingGoal)
		inputGoal.SetName(name)
		inputGoal.SetTarget(target)
		deadline := time.Now().AddDate(0, 0, daysFromNow)
		inputGoal.SetDeadline(time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, time.UTC))

		createdGoal, statusCode, err := requester.CreateSavingGoal(inputGoal)
		c.Equal(http.StatusCreated, statusCode)
		c.NoError(err, "creating saving goal failed")
		c.NotNil(createdGoal, "created saving goal is nil")
		c.NotEmpty(createdGoal.GetSavingGoalID(), "created saving goal id is empty")

		goalIDs = append(goalIDs, createdGoal.GetSavingGoalID())
		return createdGoal
	}

	createdGoals = append(createdGoals, createGoal("Goal 1", 1000, 30))   // 30 days, $1000
	createdGoals = append(createdGoals, createGoal("Goal 2", 5000, 90))   // 90 days, $5000
	createdGoals = append(createdGoals, createGoal("Goal 3", 2000, 60))   // 60 days, $2000
	createdGoals = append(createdGoals, createGoal("Goal 4", 10000, 365)) // 365 days, $10000
	createdGoals = append(createdGoals, createGoal("Goal 5", 500, 7))     // 7 days, $500

	defer func() {
		for _, id := range goalIDs {
			statusCode, err := requester.DeleteSavingGoal(id)
			if statusCode != http.StatusNoContent || err != nil {
				t.Logf("Failed to delete goal %s: %v", id, err)
			}
		}
	}()

	t.Run("Get all goals with default parameters", func(t *testing.T) {
		goals, statusCode, nextKey, err := requester.GetSavingGoals("", "", "", 10)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")
		c.Empty(nextKey, "expected no next key with page size 10")
	})

	t.Run("Sort by deadline ascending", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("deadline", "asc", "", 10)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with sorting failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")

		for i := 0; i < len(goals)-1; i++ {
			if i+1 < len(goals) {
				c.LessOrEqual(goals[i].GetDeadline().Unix(), goals[i+1].GetDeadline().Unix(),
					"goals not sorted by deadline ascending")
			}
		}
	})

	t.Run("Sort by deadline descending", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("deadline", "desc", "", 10)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with sorting failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")

		for i := 0; i < len(goals)-1; i++ {
			if i+1 < len(goals) {
				c.GreaterOrEqual(goals[i].GetDeadline().Unix(), goals[i+1].GetDeadline().Unix(),
					"goals not sorted by deadline descending")
			}
		}
	})

	t.Run("Sort by target ascending", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("target", "asc", "", 10)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with sorting failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")

		for i := 0; i < len(goals)-1; i++ {
			if i+1 < len(goals) {
				c.LessOrEqual(goals[i].GetTarget(), goals[i+1].GetTarget(),
					"goals not sorted by target ascending")
			}
		}
	})

	t.Run("Sort by target descending", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("target", "desc", "", 10)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with sorting failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")

		for i := 0; i < len(goals)-1; i++ {
			if i+1 < len(goals) {
				c.GreaterOrEqual(goals[i].GetTarget(), goals[i+1].GetTarget(),
					"goals not sorted by target descending")
			}
		}
	})

	t.Run("Invalid sort parameter", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("invalid_param", "asc", "", 10)
		c.NotEqual(http.StatusOK, statusCode)
		c.Error(err, "expected error for invalid sort parameter")
		c.Nil(goals, "goals should be nil with invalid sort parameter")
	})

	t.Run("Invalid sort order", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("deadline", "invalid_order", "", 10)
		c.NotEqual(http.StatusOK, statusCode)
		c.Error(err, "expected error for invalid sort order")
		c.Nil(goals, "goals should be nil with invalid sort order")
	})

	t.Run("Pagination with page size of 2", func(t *testing.T) {
		firstPageGoals, statusCode, nextKey, err := requester.GetSavingGoals("", "", "", 2)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get first page failed")
		c.Equal(2, len(firstPageGoals), "expected exactly 2 goals on first page")
		c.NotEmpty(nextKey, "expected next key for page size 2")

		secondPageGoals, statusCode, nextKey2, err := requester.GetSavingGoals("", "", nextKey, 2)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get second page failed")
		c.Equal(2, len(secondPageGoals), "expected exactly 2 goals on second page")
		c.NotEmpty(nextKey2, "expected next key for second page")

		for _, firstPageGoal := range firstPageGoals {
			for _, secondPageGoal := range secondPageGoals {
				c.NotEqual(firstPageGoal.GetSavingGoalID(), secondPageGoal.GetSavingGoalID(),
					"found same goal on different pages")
			}
		}

		thirdPageGoals, statusCode, nextKey3, err := requester.GetSavingGoals("", "", nextKey2, 2)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get third page failed")
		c.GreaterOrEqual(len(thirdPageGoals), 1, "expected at least 1 goal on third page")

		// If we have exactly 5 goals, we should have 1 goal on third page and no next key
		if len(thirdPageGoals) == 1 {
			c.Empty(nextKey3, "expected no next key for last page")
		}
	})

	t.Run("Small page size", func(t *testing.T) {
		goals, statusCode, nextKey, err := requester.GetSavingGoals("", "", "", 1)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with small page size failed")
		c.Equal(1, len(goals), "expected exactly 1 goal with page size 1")
		c.NotEmpty(nextKey, "expected next key with page size 1")
	})

	t.Run("Large page size", func(t *testing.T) {
		goals, statusCode, nextKey, err := requester.GetSavingGoals("", "", "", 100)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with large page size failed")
		c.GreaterOrEqual(len(goals), 5, "expected at least 5 goals")
		c.Empty(nextKey, "expected no next key with large page size")
	})

	t.Run("Invalid page size", func(t *testing.T) {
		goals, statusCode, _, err := requester.GetSavingGoals("", "", "", -1)
		c.NotEqual(http.StatusOK, statusCode)
		c.Error(err, "expected error for invalid page size")
		c.Nil(goals, "goals should be nil with invalid page size")
	})

	t.Run("All parameters combined", func(t *testing.T) {
		goals, statusCode, nextKey, err := requester.GetSavingGoals("deadline", "desc", "", 2)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get saving goals with combined parameters failed")
		c.Equal(2, len(goals), "expected exactly 2 goals")
		c.NotEmpty(nextKey, "expected next key with page size 2")

		if len(goals) > 1 {
			c.GreaterOrEqual(goals[0].GetDeadline().Unix(), goals[1].GetDeadline().Unix(),
				"goals not sorted by deadline descending")
		}

		nextPageGoals, statusCode, _, err := requester.GetSavingGoals("deadline", "desc", nextKey, 2)
		c.Equal(http.StatusOK, statusCode)
		c.NoError(err, "get next page failed")

		if len(nextPageGoals) > 0 && len(goals) > 0 {
			c.GreaterOrEqual(goals[len(goals)-1].GetDeadline().Unix(), nextPageGoals[0].GetDeadline().Unix(),
				"pagination broke sorting continuity")
		}
	})
}
