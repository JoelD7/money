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
		logger.Info(fmt.Sprintf("test_event_emitted_%d", i+1), []models.LoggerObject{user})
	}
}
