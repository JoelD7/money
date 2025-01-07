package logger

import (
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testUser struct {
	Name       string
	Age        int
	Income     float64
	DayOfBirth *time.Time
}

func (u *testUser) Key() string {
	return "user"
}

func (u *testUser) Value() map[string]interface{} {
	return map[string]interface{}{
		"s_name":         u.Name,
		"i_age":          u.Age,
		"f_income":       u.Income,
		"t_day_of_birth": u.DayOfBirth,
	}
}

func TestInfo(t *testing.T) {
	c := require.New(t)

	dayOfBirth, err := time.Parse("January 2, 2006 at 15:04:05", "April 13, 2000 at 18:23:00")
	c.Nil(err)

	name := utils.GenerateDynamoID("")

	user := &testUser{
		Name:       name,
		Age:        22,
		Income:     123456,
		DayOfBirth: &dayOfBirth,
	}

	//logger := NewLoggerMock(nil)
	logger := NewLogger()
	defer func() {
		err = logger.Finish()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println(name)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 1000)
		logger.Info(fmt.Sprintf("test_event_emitted_%d", i+1), models.Any("user", user))
	}
}

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
