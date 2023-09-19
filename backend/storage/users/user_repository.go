package users

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateUser(ctx context.Context, fullName, username, password string) error
	GetUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}
