package models

import (
	"fmt"
	"strings"
	"time"
)

type Period struct {
	Username    string     `json:"username,omitempty"`
	ID          string     `json:"period,omitempty"`
	Name        *string    `json:"name,omitempty"`
	StartDate   PeriodTime `json:"start_date,omitempty"`
	EndDate     PeriodTime `json:"end_date,omitempty"`
	CreatedDate time.Time  `json:"created_date,omitempty"`
	UpdatedDate time.Time  `json:"updated_date,omitempty"`
}

type PeriodTime struct {
	time.Time
}

const ctLayout = "2006-01-02"

func ToTime(p PeriodTime) time.Time {
	return p.Time
}

func ToPeriodTime(t time.Time) PeriodTime {
	return PeriodTime{t}
}

func (ct *PeriodTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}

	var err error

	ct.Time, err = time.Parse(ctLayout, s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal %s into a time.Time: %w", b, err)
	}

	return err
}

func (ct *PeriodTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

//var nilTime = (time.Time{}).UnixNano()

func (ct *PeriodTime) IsSet() bool {
	return !ct.IsZero()
}
