package entities

import "time"

type Person struct {
	PersonID             string        `json:"person_id,omitempty" dynamodbav:"person_id"`
	FullName             string        `json:"full_name,omitempty" dynamodbav:"full_name,omitempty"`
	Email                string        `json:"email,omitempty" dynamodbav:"email"`
	Password             string        `json:"-" dynamodbav:"password"`
	Categories           []*Category   `json:"categories,omitempty" dynamodbav:"categories,omitempty"`
	SavingGoals          []*SavingGoal `json:"saving_goals,omitempty" dynamodbav:"saving_goals,omitempty"`
	CreatedDate          time.Time     `json:"created_date,omitempty" dynamodbav:"created_date,omitempty"`
	UpdatedDate          time.Time     `json:"updated_date,omitempty" dynamodbav:"update_date,omitempty"`
	PreviousRefreshToken string        `json:"previous_refresh_token,omitempty" dynamodbav:"previous_refresh_token"`
	RefreshToken         string        `json:"refresh_token,omitempty" dynamodbav:"refresh_token"`
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
