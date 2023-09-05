package users

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	CreateUser(ctx context.Context, fullName, email, password string) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}
