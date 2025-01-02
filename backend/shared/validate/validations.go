package validate

import (
	"github.com/JoelD7/money/backend/models"
	"math"
	"regexp"
)

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"

var (
	validSortBy = map[string]struct{}{
		string(models.SortParamCreatedDate): {},
		string(models.SortParamAmount):      {},
		string(models.SortParamName):        {},
	}
)

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

func SortBy(sortBy string) error {
	if _, ok := validSortBy[sortBy]; !ok && sortBy != "" {
		return models.ErrInvalidSortBy
	}

	return nil
}

func SortType(sortType string) error {
	if sortType == string(models.SortOrderDescending) || sortType == string(models.SortOrderAscending) || sortType == "" {
		return nil
	}

	return models.ErrInvalidSortOrder
}
