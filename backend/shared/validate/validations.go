package validate

import (
	"github.com/JoelD7/money/backend/models"
	"math"
	"regexp"
)

const (
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"

	SortByModelExpenses    SortByModel = "expenses"
	SortByModelSavingGoals SortByModel = "saving_goals"
)

type SortByModel string

var (
	validSortBy = map[SortByModel]map[string]struct{}{
		SortByModelExpenses: {
			string(models.SortParamCreatedDate): {},
			string(models.SortParamAmount):      {},
			string(models.SortParamName):        {},
		},
		SortByModelSavingGoals: {
			string(models.SortParamDeadline): {},
			string(models.SortParamTarget):   {},
		},
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

func SortBy(sortBy string, model SortByModel) error {
	if _, ok := validSortBy[model][sortBy]; !ok && sortBy != "" {
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

func PageSize(pageSize int) error {
	if pageSize < 0 || pageSize > math.MaxInt32 {
		return models.ErrInvalidPageSize
	}

	return nil
}
