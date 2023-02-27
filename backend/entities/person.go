package entities

import "time"

type Person struct {
	PersonID    string        `json:"person_id,omitempty" dynamodbav:"person_id"`
	FullName    string        `json:"full_name,omitempty" dynamodbav:"full_name"`
	Email       string        `json:"email,omitempty" dynamodbav:"email"`
	Password    string        `json:"-" dynamodbav:"password"`
	Categories  []*Category   `json:"categories,omitempty" dynamodbav:"categories"`
	SavingGoals []*SavingGoal `json:"saving_goals,omitempty" dynamodbav:"saving_goals"`
	CreatedDate time.Time     `json:"created_date,omitempty" dynamodbav:"created_date"`
	UpdatedDate time.Time     `json:"updated_date,omitempty" dynamodbav:"update_date"`
}

type Category struct {
	CategoryID    string         `json:"category_id,omitempty" dynamodbav:"category_id"`
	CategoryName  string         `json:"category_name,omitempty" dynamodbav:"category_name"`
	Budget        float64        `json:"budget,omitempty" dynamodbav:"budget"`
	Color         string         `json:"color,omitempty" dynamodbav:"color"`
	Keywords      []string       `json:"keywords,omitempty" dynamodbav:"keywords,stringset"`
	Subcategories []*Subcategory `json:"subcategories,omitempty" dynamodbav:"subcategories"`
}

type Subcategory struct {
	SubcategoryID   string `json:"subcategory,omitempty" dynamodbav:"subcategory"`
	SubcategoryName string `json:"subcategory_name,omitempty" dynamodbav:"subcategory_name"`
	Color           string `json:"color,omitempty" dynamodbav:"color"`
}

type SavingGoal struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty" dynamodbav:"saving_goal_id"`
	Name         string    `json:"name,omitempty" dynamodbav:"name"`
	Goal         float64   `json:"goal,omitempty" dynamodbav:"goal"`
	Deadline     time.Time `json:"deadline,omitempty" dynamodbav:"deadline"`
}
