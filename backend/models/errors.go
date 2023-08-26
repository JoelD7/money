package models

import "errors"

var (
	ErrMissingEmail         = errors.New("missing email")
	ErrMissingPassword      = errors.New("missing password")
	ErrInvalidEmail         = errors.New("email is invalid")
	ErrUserNotFound         = errors.New("user not found")
	ErrExistingUser         = errors.New("this account already exists")
	ErrWrongCredentials     = errors.New("the email or password are incorrect")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRefreshTokenMismatch = errors.New("received refresh token doesn't match with the user's")
	ErrUnauthorized         = errors.New("Unauthorized")
	ErrSigningKeyNotFound   = errors.New("signing key not found")
	ErrTokensNotFound       = errors.New("no invalid tokens found")
)
