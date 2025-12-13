package users

import (
	"testing"

	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestGetUpdateParams(t *testing.T) {
	t.Run("User with CurrentPeriod set", func(t *testing.T) {
		currentPeriod := "2024-01"
		user := &models.User{
			CurrentPeriod: &currentPeriod,
		}

		expectedVal, err := attributevalue.Marshal(currentPeriod)
		assert.NoError(t, err)

		expression, values, err := getUpdateParams(user)

		assert.NoError(t, err)
		assert.Equal(t, "SET :current_period", expression)

		assert.Len(t, values, 1)
		assert.Contains(t, values, ":current_period")
		assert.Equal(t, expectedVal, values[":current_period"])
	})

	t.Run("User with nil CurrentPeriod", func(t *testing.T) {
		user := &models.User{
			CurrentPeriod: nil,
		}

		expression, values, err := getUpdateParams(user)
		assert.ErrorContains(t, err, "no attributes to update")
		assert.Empty(t, expression)
		assert.Empty(t, values)
	})

	t.Run("User with empty string CurrentPeriod", func(t *testing.T) {
		currentPeriod := ""
		user := &models.User{
			CurrentPeriod: &currentPeriod,
		}

		expectedVal := &types.AttributeValueMemberS{Value: ""}

		expression, values, err := getUpdateParams(user)

		assert.NoError(t, err)
		assert.Equal(t, "SET :current_period", expression)
		assert.Len(t, values, 1)
		assert.Equal(t, expectedVal, values[":current_period"])
	})
}
