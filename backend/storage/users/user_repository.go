package users

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type RepositoryAPI interface {
	createUser(ctx context.Context, fullName, email, password string) error
	getUser(ctx context.Context, userID string) (*models.User, error)
	getUserByEmail(ctx context.Context, email string) (*models.User, error)
	updateUser(ctx context.Context, user *models.User) error
}

type Repository struct {
	client RepositoryAPI
}

func NewRepository(client RepositoryAPI) *Repository {
	return &Repository{client}
}

func (u *Repository) CreateUser(ctx context.Context, fullName, email, password string) error {
	return u.client.createUser(ctx, fullName, email, password)
}

func (u *Repository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return u.client.getUser(ctx, userID)
}

func (u *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.client.getUserByEmail(ctx, email)
}

func (u *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	return u.client.updateUser(ctx, user)
}
