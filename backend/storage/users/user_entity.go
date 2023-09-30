package users

import (
	"github.com/JoelD7/money/backend/models"
	"time"
)

type userEntity struct {
	FullName      string            `json:"full_name,omitempty" dynamodbav:"full_name,omitempty"`
	Username      string            `json:"username,omitempty" dynamodbav:"username"`
	Password      string            `json:"-" dynamodbav:"password"`
	Categories    []*categoryEntity `json:"categories,omitempty" dynamodbav:"categories,omitempty"`
	CreatedDate   time.Time         `json:"created_date,omitempty" dynamodbav:"created_date,omitempty"`
	UpdatedDate   time.Time         `json:"updated_date,omitempty" dynamodbav:"update_date,omitempty"`
	AccessToken   string            `json:"-" dynamodbav:"access_token,omitempty"`
	RefreshToken  string            `json:"-" dynamodbav:"refresh_token"`
	CurrentPeriod string            `json:"current_period,omitempty" dynamodbav:"current_period,omitempty"`
}

type categoryEntity struct {
	ID     string   `json:"id,omitempty" dynamodbav:"id"`
	Name   *string  `json:"name,omitempty" dynamodbav:"name"`
	Budget *float64 `json:"budget,omitempty" dynamodbav:"budget,omitempty"`
	Color  *string  `json:"color,omitempty" dynamodbav:"color,omitempty"`
}

func toUserEntity(u *models.User) *userEntity {
	return &userEntity{
		FullName:      u.FullName,
		Username:      u.Username,
		Password:      u.Password,
		Categories:    toCategoryEntities(u.Categories),
		CreatedDate:   u.CreatedDate,
		UpdatedDate:   u.UpdatedDate,
		AccessToken:   u.AccessToken,
		RefreshToken:  u.RefreshToken,
		CurrentPeriod: u.CurrentPeriod,
	}
}

func toCategoryEntities(modelCategories []*models.Category) []*categoryEntity {
	categories := make([]*categoryEntity, len(modelCategories))

	for _, v := range modelCategories {
		categories = append(categories, toCategoryEntity(v))
	}

	return categories
}

func toCategoryEntity(modelCategory *models.Category) *categoryEntity {
	return &categoryEntity{
		ID:     modelCategory.ID,
		Name:   modelCategory.Name,
		Budget: modelCategory.Budget,
		Color:  modelCategory.Color,
	}
}

func toUserModel(u *userEntity) *models.User {
	return &models.User{
		FullName:      u.FullName,
		Username:      u.Username,
		Password:      u.Password,
		Categories:    toCategoryModels(u.Categories),
		CreatedDate:   u.CreatedDate,
		UpdatedDate:   u.UpdatedDate,
		AccessToken:   u.AccessToken,
		RefreshToken:  u.RefreshToken,
		CurrentPeriod: u.CurrentPeriod,
	}
}

func toCategoryModels(entityCategories []*categoryEntity) []*models.Category {
	categories := make([]*models.Category, len(entityCategories))

	for _, v := range entityCategories {
		categories = append(categories, toCategoryModel(v))
	}

	return categories
}

func toCategoryModel(entityCategory *categoryEntity) *models.Category {
	return &models.Category{
		ID:     entityCategory.ID,
		Name:   entityCategory.Name,
		Budget: entityCategory.Budget,
		Color:  entityCategory.Color,
	}
}
