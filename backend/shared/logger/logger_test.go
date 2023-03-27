package logger

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testPerson struct {
	Name       string
	Age        int
	Income     float64
	DayOfBirth *time.Time
}

func (p *testPerson) LogName() string {
	return "person"
}

func (p *testPerson) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_name":         p.Name,
		"i_age":          p.Age,
		"f_income":       p.Income,
		"t_day_of_birth": p.DayOfBirth,
	}
}

func TestInfo(t *testing.T) {
	c := require.New(t)

	dayOfBirth, err := time.Parse("January 2, 2006 at 15:04:05", "April 13, 2000 at 18:23:00")
	c.Nil(err)

	person := &testPerson{
		Name:       "Joel",
		Age:        22,
		Income:     123456,
		DayOfBirth: &dayOfBirth,
	}

	logger := NewLogger()

	logger.Info("test_event_emitted", []Object{person})
}
