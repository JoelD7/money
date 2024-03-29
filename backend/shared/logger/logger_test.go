package logger

import (
	"github.com/JoelD7/money/backend/models"
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

func (u *testUser) LogName() string {
	return "user"
}

func (u *testUser) LogProperties() map[string]interface{} {
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

	user := &testUser{
		Name:       "Joel",
		Age:        22,
		Income:     123456,
		DayOfBirth: &dayOfBirth,
	}

	logger := NewLoggerMock(nil)
	defer func() {
		err = logger.Close()
		if err != nil {
			panic(err)
		}
	}()

	logger.Info("test_event_emitted", []models.LoggerObject{user})
}
