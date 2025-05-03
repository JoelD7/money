package models

import (
	"errors"
)

var (
	// User
	ErrMissingUsername       = errors.New("missing username")
	ErrMissingPassword       = errors.New("missing password")
	ErrInvalidEmail          = errors.New("invalid username. username must be a valid email address")
	ErrUserNotFound          = errors.New("user not found")
	ErrExistingUser          = errors.New("this account already exists")
	ErrWrongCredentials      = errors.New("the email or password are incorrect")
	ErrInvalidToken          = errors.New("invalid token")
	ErrRefreshTokenMismatch  = errors.New("received refresh token doesn't match with the user's")
	ErrUnauthorized          = errors.New("Unauthorized")
	ErrSigningKeyNotFound    = errors.New("signing key not found")
	ErrInvalidTokensNotFound = errors.New("no invalid tokens found")
	ErrSecretNotFound        = errors.New("secret not found")
	// ErrMalformedToken error when the client sends a token that doesn't comply with the JWT standard.
	// This message is included for security reasons. We aim to give the client minimal information about why the request
	// was denied. If we were to state that 'this token is malformed,' it could signal an attacker that the denial was
	// linked to the token's structure or content, inadvertently assisting him.
	ErrMalformedToken = errors.New("invalid token")
	// ErrUsernameDeleteMismatch error when the username in the authorization header doesn't match with the username in
	// the path parameter of an endpoint like /users/{username}. This error is currently returned when a user tries to
	// delete another user.
	ErrUsernameDeleteMismatch = errors.New("authorization username doesn't match with path parameter username")

	// Income
	ErrIncomeNotFound        = errors.New("user income not found")
	ErrExistingIncome        = errors.New("this income already exists")
	ErrMissingIncomeID       = errors.New("missing income id")
	ErrIncomePeriodsNotFound = errors.New("income periods not found")

	// General
	ErrInvalidAmount            = errors.New("invalid amount. The amount has to be a number greater than zero")
	ErrInvalidRequestBody       = errors.New("invalid request body")
	ErrMissingAmount            = errors.New("missing amount")
	ErrInvalidPageSize          = errors.New("invalid page size")
	ErrInvalidStartKey          = errors.New("invalid start key")
	ErrNoUsernameInContext      = errors.New("username not found in authorizer context")
	ErrMissingName              = errors.New("missing name")
	ErrInvalidSortOrder         = errors.New("invalid sort order")
	ErrInvalidSortBy            = errors.New("invalid sort by")
	ErrNoMoreItemsToBeRetrieved = errors.New("no more items to be retrieved")
	// ErrIndexKeysNotFound error when a DynamoDB index is not included in the map used to build the LastEvaluatedKey.
	// This is important as it breaks pagination.
	ErrIndexKeysNotFound     = errors.New("index keys not found")
	ErrMissingIdempotencyKey = errors.New("missing idempotency key header")
	// ErrUnexpectedTypeAssertion error when the type of a resource obtained from the idempotency manager is unexpected.
	ErrUnexpectedTypeAssertion = errors.New("unexpected type assertion in idempotency manager")

	// Saving
	ErrMissingSavingID      = errors.New("missing saving id")
	ErrInvalidSavingAmount  = errors.New("saving amount must be greater than or equal to 0")
	ErrSavingsNotFound      = errors.New("user savings not found")
	ErrSavingNotFound       = errors.New("saving not found")
	ErrUpdateSavingNotFound = errors.New("the saving you are trying to update does not exist")
	ErrDeleteSavingNotFound = errors.New("the saving you are trying to delete does not exist")

	// Category
	ErrCategoriesNotFound        = errors.New("categories not found")
	ErrCategoryNotFound          = errors.New("category not found")
	ErrMissingCategoryName       = errors.New("missing category name")
	ErrMissingCategoryColor      = errors.New("missing color")
	ErrInvalidHexColor           = errors.New("invalid hex color")
	ErrInvalidBudget             = errors.New("budget must be greater than or equal to 0")
	ErrMissingCategoryBudget     = errors.New("missing budget")
	ErrCategoryNameAlreadyExists = errors.New("category name already exists")

	// Saving Goal
	ErrSavingGoalNameSettingFailed      = errors.New("saving goal name not set")
	ErrSavingGoalNotFound               = errors.New("saving goal not found")
	ErrSavingGoalsNotFound              = errors.New("saving goals not found")
	ErrMissingSavingGoalName            = errors.New("missing goal name")
	ErrMissingSavingGoalTarget          = errors.New("missing goal target")
	ErrInvalidSavingGoalTarget          = errors.New("goal target must be greater than 0")
	ErrInvalidSavingGoalDeadline        = errors.New("deadline must be in the future")
	ErrMissingSavingGoalRecurringAmount = errors.New("missing recurring amount")

	// Expense
	ErrMissingExpenseID          = errors.New("missing expense id")
	ErrMissingRecurringDay       = errors.New("missing recurring day")
	ErrInvalidRecurringDay       = errors.New("recurring day must be between 1 and 31")
	ErrCategoryNameSettingFailed = errors.New("couldn't set category name for expenses")
	ErrRecurringExpenseNameTaken = errors.New("recurring expense name is taken")
	ErrRecurringExpensesNotFound = errors.New("recurring expenses not found")
	ErrRecurringExpenseNotFound  = errors.New("recurring expense not found")
	ErrMissingExpenseRecurringID = errors.New("missing expense recurring id")
	ErrExpenseNotFound           = errors.New("expense not found")
	ErrExpensesNotFound          = errors.New("user expenses not found")

	// Period
	ErrPeriodNotFound                 = errors.New("period not found")
	ErrPeriodsNotFound                = errors.New("periods not found")
	ErrInvalidPeriod                  = errors.New("invalid period")
	ErrMissingPeriodDates             = errors.New("missing period dates. A period should have a start_date and end_date")
	ErrStartDateShouldBeBeforeEndDate = errors.New("start_date should be before end_date")
	ErrPeriodNameIsTaken              = errors.New("period name is taken")
	ErrUpdatePeriodNotFound           = errors.New("the period you are trying to update does not exist")
	ErrInvalidPeriodDate              = errors.New("invalid period date")
	ErrMissingPeriodID                = errors.New("missing period id")
	ErrMissingPeriodName              = errors.New("missing period name")
	ErrMissingPeriodStartDate         = errors.New("missing period start date")
	ErrMissingPeriod                  = errors.New("missing period")
	ErrMissingPeriodCreatedDate       = errors.New("missing period created date")
	ErrMissingPeriodUpdatedDate       = errors.New("missing period updated date")
)
