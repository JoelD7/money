package models

import (
	"errors"
)

var (
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
	ErrMalformedToken              = errors.New("invalid token")
	ErrSavingsNotFound             = errors.New("user savings not found")
	ErrSavingNotFound              = errors.New("saving not found")
	ErrUpdateSavingNotFound        = errors.New("the saving you are trying to update does not exist")
	ErrDeleteSavingNotFound        = errors.New("the saving you are trying to delete does not exist")
	ErrIncomeNotFound              = errors.New("user income not found")
	ErrMissingPeriod               = errors.New("missing period")
	ErrExpenseNotFound             = errors.New("expense not found")
	ErrExpensesNotFound            = errors.New("user expenses not found")
	ErrInvalidAmount               = errors.New("invalid amount. The amount has to be a number greater than zero")
	ErrInvalidRequestBody          = errors.New("invalid request body")
	ErrMissingSavingID             = errors.New("missing saving id")
	ErrInvalidPageSize             = errors.New("invalid page size")
	ErrInvalidStartKey             = errors.New("invalid start key")
	ErrCategoriesNotFound          = errors.New("categories not found")
	ErrCategoryNotFound            = errors.New("category not found")
	ErrMissingCategoryName         = errors.New("name should not be empty")
	ErrMissingCategoryColor        = errors.New("missing color")
	ErrInvalidHexColor             = errors.New("invalid hex color")
	ErrInvalidBudget               = errors.New("budget must be greater than or equal to 0")
	ErrMissingCategoryBudget       = errors.New("missing budget")
	ErrCategoryNameAlreadyExists   = errors.New("category name already exists")
	ErrMissingAmount               = errors.New("missing amount")
	ErrInvalidSavingAmount         = errors.New("saving amount must be greater than or equal to 0")
	ErrSavingGoalNameSettingFailed = errors.New("saving goal name not set")
	ErrSavingGoalNotFound          = errors.New("saving goal not found")
	ErrNoUsernameInContext         = errors.New("couldn't identify the user. Check if your Bearer token header is correct")
	ErrCategoryNameSettingFailed   = errors.New("category name not set")
	ErrMissingName                 = errors.New("missing name")
	ErrMissingExpenseID            = errors.New("missing expense id")

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
	ErrMissingPeriodCreatedDate       = errors.New("missing period created date")
	ErrMissingPeriodUpdatedDate       = errors.New("missing period updated date")
)
