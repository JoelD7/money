package main

import (
	"github.com/JoelD7/money/backend/models"
	"math"
	"regexp"
)

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"

func Email(email string) error {
	regex := regexp.MustCompile(emailRegex)

	if email == "" {
		return models.ErrMissingUsername
	}

	if !regex.MatchString(email) {
		return models.ErrInvalidEmail
	}

	return nil
}

func Amount(amount *float64) error {
	if amount != nil && (*amount <= 0 || *amount > math.MaxFloat64) {
		return models.ErrInvalidAmount
	}

	return nil
}
