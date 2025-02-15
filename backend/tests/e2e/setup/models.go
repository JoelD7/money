package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/storage/users"
	"os"
)

const (
	// Username is the value to be used as "username" for all e2e tests
	Username = "e2e_test@gmail.com"

	samplesDir = "/Users/joelfabian/go/src/github.com/JoelD7/money/backend/tests/e2e/setup/samples"
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

	_, err = repo.CreateUser(ctx, user)
	if err == nil {
		userCreated = true
	}

	return
}

// CreateIncomeEntries creates income entries in the DB taken from "source" file. Deletes them after tests are completed
func CreateIncomeEntries(ctx context.Context, repo income.Repository, source string, cleaner Cleaner) (entries []*models.Income, err error) {
	incomeCreated := false

	if source == "" {
		source = samplesDir + "/income.json"
	}

	defer cleaner.Cleanup(func() {
		if !incomeCreated {
			return
		}

		err = repo.BatchDeleteIncome(ctx, entries)
	})

	data, err := os.ReadFile(source)
	if err != nil {
		err = fmt.Errorf("cannot create test income entries: %v", err)
		return
	}

	err = json.Unmarshal(data, &entries)
	if err != nil {
		err = fmt.Errorf("cannot create test income entries: %v", err)
		return
	}

	err = repo.BatchCreateIncome(ctx, entries)
	if err != nil {
		err = fmt.Errorf("cannot create test income entries: %v", err)
		return
	}

	incomeCreated = true
	return
}

func CreateExpensesEntries(ctx context.Context, repo expenses.Repository, source string, cleaner Cleaner) (entries []*models.Expense, err error) {
	expensesCreated := false

	if source == "" {
		source = samplesDir + "/expenses.json"
	}

	defer cleaner.Cleanup(func() {
		if !expensesCreated {
			return
		}

		err = repo.BatchDeleteExpenses(ctx, entries)
	})

	data, err := os.ReadFile(source)
	if err != nil {
		err = fmt.Errorf("cannot create test expenses entries: %v", err)
		return
	}

	err = json.Unmarshal(data, &entries)
	if err != nil {
		err = fmt.Errorf("cannot create test expenses entries: %v", err)
		return
	}

	err = repo.BatchCreateExpenses(ctx, entries)
	if err != nil {
		err = fmt.Errorf("cannot create test expenses entries: %v", err)
		return
	}

	expensesCreated = true
	return
}
