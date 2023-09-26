package models

import "errors"

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
	ErrMalformedToken       = errors.New("invalid token")
	ErrSavingsNotFound      = errors.New("savings not found")
	ErrUpdateSavingNotFound = errors.New("the saving you are trying to update does not exist")
	ErrDeleteSavingNotFound = errors.New("the saving you are trying to delete does not exist")
	ErrIncomeNotFound       = errors.New("user income not found")
	ErrMissingPeriod        = errors.New("missing period")
	ErrExpensesNotFound     = errors.New("user expenses not found")
	ErrInvalidAmount        = errors.New("invalid amount. The amount has to be a number greater than zero")
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrMissingSavingID      = errors.New("missing saving id")
	ErrInvalidPageSize      = errors.New("invalid page size")
	ErrInvalidStartKey      = errors.New("invalid start key")
	ErrCategoriesNotFound   = errors.New("categories not found")
)
