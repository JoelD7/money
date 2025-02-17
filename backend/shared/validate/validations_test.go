package validate

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortBy_Expenses(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		model    SortByModel
		expected error
	}{
		{
			name:     "Valid sort by - created_date",
			input:    "created_date",
			model:    SortByModelExpenses,
			expected: nil,
		},
		{
			name:     "Valid sort by - amount",
			input:    "amount",
			model:    SortByModelExpenses,
			expected: nil,
		},
		{
			name:     "Valid sort by - name",
			input:    "name",
			model:    SortByModelExpenses,
			expected: nil,
		},
		{
			name:     "Invalid sort by - random value",
			input:    "random",
			model:    SortByModelExpenses,
			expected: models.ErrInvalidSortBy,
		},
		{
			name:     "Invalid sort by - empty string",
			input:    "",
			model:    SortByModelExpenses,
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SortBy(tc.input, tc.model)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestSortBy_SavingGoals(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		model    SortByModel
		expected error
	}{
		{
			name:     "Valid sort by - created_date",
			input:    "created_date",
			model:    SortByModelSavingGoals,
			expected: nil,
		},
		{
			name:     "Valid sort by - amount",
			input:    "amount",
			model:    SortByModelSavingGoals,
			expected: nil,
		},
		{
			name:     "Invalid sort by - name",
			input:    "name",
			model:    SortByModelSavingGoals,
			expected: models.ErrInvalidSortBy,
		},
		{
			name:     "Invalid sort by - random value",
			input:    "random",
			model:    SortByModelSavingGoals,
			expected: models.ErrInvalidSortBy,
		},
		{
			name:     "Invalid sort by - empty string",
			input:    "",
			model:    SortByModelSavingGoals,
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SortBy(tc.input, tc.model)
			assert.Equal(t, tc.expected, err)
		})
	}
}
