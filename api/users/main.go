package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"time"
)

type Person struct {
	PersonID    string        `json:"person_id,omitempty"`
	FullName    string        `json:"full_name,omitempty"`
	Email       string        `json:"email,omitempty"`
	Categories  []*Category   `json:"categories,omitempty"`
	SavingGoals []*SavingGoal `json:"saving_goals,omitempty"`
}

type Category struct {
	CategoryID    string         `json:"category_id,omitempty"`
	CategoryName  string         `json:"category_name,omitempty"`
	Budget        float64        `json:"budget,omitempty"`
	Color         string         `json:"color,omitempty"`
	Keywords      []string       `json:"keywords,omitempty" dynamodbav:"keywords,stringset"`
	Subcategories []*Subcategory `json:"subcategories,omitempty"`
}

type Subcategory struct {
	SubcategoryID   string `json:"subcategory,omitempty"`
	SubcategoryName string `json:"subcategory_name,omitempty"`
	Color           string `json:"color,omitempty"`
}

type SavingGoal struct {
	SavingGoalID string    `json:"saving_goal_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Goal         float64   `json:"goal,omitempty"`
	Deadline     time.Time `json:"deadline,omitempty"`
}

func dummy() (*Person, error) {
	person, err := getItem("personId")
	if err != nil {
		return nil, err
	}

	return person, nil
}

func main() {
	lambda.Start(dummy)
}
