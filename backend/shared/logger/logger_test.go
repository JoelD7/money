package logger

import (
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetLogDataAsBytes(t *testing.T) {
	c := require.New(t)

	l := &Log{
		Service: "unit-test",
	}

	var data []byte
	fields := make([]models.LoggerField, 1)

	t.Run("String field", func(t *testing.T) {
		fields[0] = models.Any("key", "value")
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, fields)
		c.NotNil(data)
		c.Contains(string(data), `"properties":{"key":"value"}`)
	})

	t.Run("Struct field - implemented interface", func(t *testing.T) {
		user := models.User{
			FullName:      "John Doe",
			Username:      "johndoe123",
			Password:      "securepassword",
			CreatedDate:   time.Date(2025, 1, 5, 20, 40, 56, 0, time.UTC),
			UpdatedDate:   time.Date(2025, 1, 5, 20, 40, 56, 0, time.UTC),
			AccessToken:   "access-token-placeholder",
			RefreshToken:  "refresh-token-placeholder",
			CurrentPeriod: "2025-01",
			Remainder:     1500.0,
		}

		fields[0] = models.Any("user", user)
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, fields)
		c.NotNil(data)
		c.Contains(string(data), `"properties":{"user":{"full_name":"John Doe","username":"johndoe123","created_date":"2025-01-05T20:40:56Z","updated_date":"2025-01-05T20:40:56Z","current_period":"2025-01","remainder":1500}`)
	})

	t.Run("Map field", func(t *testing.T) {
		fields[0] = models.Any("key", map[string]interface{}{"subkey": "value"})
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, fields)
		c.NotNil(data)
		c.Contains(string(data), `"properties":{"key":{"subkey":"value"}}`)
	})

	t.Run("Array field", func(t *testing.T) {
		fields[0] = models.Any("key", []string{"value1", "value2"})
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, fields)
		c.NotNil(data)
		c.Contains(string(data), `"properties":{"key":["value1","value2"]}`)
	})

	t.Run("Nil field", func(t *testing.T) {
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, []models.LoggerField{nil})
		c.NotNil(data)
		c.Contains(string(data), `"event":"test_event"`)
	})

	t.Run("Struct field", func(t *testing.T) {
		type authRequestBody struct {
			Username string
			Password string
		}

		fields[0] = models.Any("request_body", authRequestBody{"username", "password"})
		data = l.getLogDataAsBytes(infoLevel, "test_event", nil, fields)
		c.NotNil(data)
		fmt.Println(string(data))
	})
}
