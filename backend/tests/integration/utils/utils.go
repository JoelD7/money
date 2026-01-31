package utils

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/users"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	//minBound and maxBound define the inclusive range for generating random username suffixes in tests.
	minBound = 1000
	maxBound = 2000

	usernamePrefix = "integration_test_user"
	fullnamePrefix = "Integration Test User"
)

func CreateUser(ctx context.Context, dynamoClient *dynamodb.Client, envConfig *models.EnvironmentConfiguration, t *testing.T) (createdUser *models.User, err error) {
	var repo *users.DynamoRepository

	defer func() {
		//If repo is nil, so will createdUser. Added nil check for "repo" to make this dependency explicit and not
		//breaking the code accidentally.
		if createdUser == nil || repo == nil {
			return
		}

		t.Cleanup(func() {
			deleteUserErr := repo.DeleteUser(ctx, createdUser.Username)
			if deleteUserErr != nil {
				err = fmt.Errorf("failed to delete created user during cleanup: %w", deleteUserErr)
			}
		})
	}()

	repo, err = users.NewDynamoRepository(dynamoClient, envConfig.UsersTable)
	if err != nil {
		return nil, err
	}

	username, fullName := getRandomUserData()

	user := &models.User{
		FullName:    fullName,
		Username:    username,
		CreatedDate: time.Now(),
	}

	createdUser, err = repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("creating testing user: %w", err)
	}

	return createdUser, nil
}

func getRandomUserData() (string, string) {
	n := rand.Intn(maxBound-minBound+1) + minBound
	return fmt.Sprintf("%s_%d", usernamePrefix, n), fmt.Sprintf("%s_%d", fullnamePrefix, n)
}
