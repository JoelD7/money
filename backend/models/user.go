package models

import "time"

type User struct {
	FullName      string      `json:"full_name,omitempty" dynamodbav:"full_name,omitempty"`
	Username      string      `json:"username,omitempty" dynamodbav:"username"`
	Password      string      `json:"-" dynamodbav:"password"`
	Categories    []*Category `json:"categories,omitempty" dynamodbav:"categories,omitempty"`
	CreatedDate   time.Time   `json:"created_date,omitempty" dynamodbav:"created_date,omitempty"`
	UpdatedDate   time.Time   `json:"updated_date,omitempty" dynamodbav:"update_date,omitempty"`
	AccessToken   string      `json:"-" dynamodbav:"access_token,omitempty"`
	RefreshToken  string      `json:"-" dynamodbav:"refresh_token"`
	CurrentPeriod string      `json:"current_period,omitempty" dynamodbav:"current_period,omitempty"`
	Remainder     float64     `json:"remainder,omitempty" dynamodbav:"remainder,omitempty"`
}

type Category struct {
	ID            string         `json:"id,omitempty" dynamodbav:"id"`
	Name          string         `json:"name,omitempty" dynamodbav:"name"`
	Budget        float64        `json:"budget,omitempty" dynamodbav:"budget,omitempty"`
	Color         string         `json:"color,omitempty" dynamodbav:"color,omitempty"`
	Keywords      []string       `json:"keywords,omitempty" dynamodbav:"keywords,stringset,omitempty"`
	Subcategories []*Subcategory `json:"subcategories,omitempty" dynamodbav:"subcategories,omitempty"`
}

type Subcategory struct {
	ID    string `json:"id,omitempty" dynamodbav:"id"`
	Name  string `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Color string `json:"color,omitempty" dynamodbav:"color,omitempty"`
}

func (u User) LogName() string {
	return "user"
}

func (u User) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_username": u.Username,
		"s_fullname": u.FullName,
	}
}
