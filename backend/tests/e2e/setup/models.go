package setup

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/users"
)

const (
	// Username is the value to be used as "username" for all e2e tests
	Username = "e2e_test@gmail.com"
)

// Cleaner describes a function that cleans up resources after a test finishes
type Cleaner interface {
	Cleanup(f func())
}

// CreateUser persists a user to the DB. Deletes it after the test is completed.
func CreateUser(ctx context.Context, repo users.Repository, user *models.User, cleaner Cleaner) (err error) {
	//Enforcing this username in e2e tests for consistency
	user.Username = Username

	userCreated := false

	defer cleaner.Cleanup(func() {
		if !userCreated {
			return
		}

		err = repo.DeleteUser(ctx, user.Username)
	})

	err = repo.CreateUser(ctx, user)
	if err == nil {
		userCreated = true
	}

	return
}
