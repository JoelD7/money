package models

import "time"

type User struct {
	FullName      string        `json:"full_name,omitempty" dynamodbav:"full_name,omitempty"`
	Username      string        `json:"username,omitempty" dynamodbav:"username"`
	Password      string        `json:"-" dynamodbav:"password"`
	Categories    []*Category   `json:"categories,omitempty" dynamodbav:"categories,omitempty"`
	SavingGoals   []*SavingGoal `json:"saving_goals,omitempty" dynamodbav:"saving_goals,omitempty"`
	CreatedDate   time.Time     `json:"created_date,omitempty" dynamodbav:"created_date,omitempty"`
	UpdatedDate   time.Time     `json:"updated_date,omitempty" dynamodbav:"update_date,omitempty"`
	AccessToken   string        `json:"-" dynamodbav:"access_token,omitempty"`
	RefreshToken  string        `json:"-" dynamodbav:"refresh_token"`
	CurrentPeriod string        `json:"current_period,omitempty" dynamodbav:"current_period,omitempty"`
	Remainder     float64       `json:"remainder,omitempty" dynamodbav:"remainder,omitempty"`
}

type Category struct {
	CategoryID    string         `json:"category_id,omitempty" dynamodbav:"category_id"`
	CategoryName  string         `json:"category_name,omitempty" dynamodbav:"category_name"`
	Budget        float64        `json:"budget,omitempty" dynamodbav:"budget,omitempty"`
	Color         string         `json:"color,omitempty" dynamodbav:"color,omitempty"`
	Keywords      []string       `json:"keywords,omitempty" dynamodbav:"keywords,stringset,omitempty"`
	Subcategories []*Subcategory `json:"subcategories,omitempty" dynamodbav:"subcategories,omitempty"`
}

type Subcategory struct {
	SubcategoryID   string `json:"subcategory_id,omitempty" dynamodbav:"subcategory_id"`
	SubcategoryName string `json:"subcategory_name,omitempty" dynamodbav:"subcategory_name,omitempty"`
	Color           string `json:"color,omitempty" dynamodbav:"color,omitempty"`
}

type SavingGoal struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty" dynamodbav:"saving_goal_id,omitempty"`
	Name         string    `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Goal         float64   `json:"goal,omitempty" dynamodbav:"goal,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty" dynamodbav:"deadline,omitempty"`
}

func (u User) LogName() string {
	return "user"
}

func (u User) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_email":    u.Username,
		"s_fullname": u.FullName,
	}
}
