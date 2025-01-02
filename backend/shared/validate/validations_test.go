package validate

import (
	"github.com/JoelD7/money/backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortBy(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Valid sort by - created_date",
			input:    "created_date",
			expected: nil,
		},
		{
			name:     "Valid sort by - amount",
			input:    "amount",
			expected: nil,
		},
		{
			name:     "Valid sort by - name",
			input:    "name",
			expected: nil,
		},
		{
			name:     "Invalid sort by - random value",
			input:    "random",
			expected: models.ErrInvalidSortBy,
		},
		{
			name:     "Invalid sort by - empty string",
			input:    "",
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SortBy(tc.input)
			assert.Equal(t, tc.expected, err)
		})
	}
}
