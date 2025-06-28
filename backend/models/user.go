package models

import "time"

type User struct {
	FullName      string      `json:"full_name,omitempty"`
	Username      string      `json:"username,omitempty"`
	Password      string      `json:"-"`
	Categories    []*Category `json:"categories,omitempty"`
	CreatedDate   time.Time   `json:"created_date,omitempty"`
	UpdatedDate   time.Time   `json:"updated_date,omitempty"`
	AccessToken   string      `json:"-"`
	RefreshToken  string      `json:"-"`
	CurrentPeriod string      `json:"current_period,omitempty"`
	Remainder     float64     `json:"remainder"`
}

type Category struct {
	ID       string   `json:"id,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Budget   *float64 `json:"budget,omitempty"`
	Color    *string  `json:"color,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
}

func (u *User) GetKey() string {
	return "user"
}

func (u *User) GetValue() (interface{}, error) {
	return map[string]interface{}{
		"s_username":     u.Username,
		"s_fullname":     u.FullName,
		"t_created_date": u.CreatedDate,
	}, nil
}
