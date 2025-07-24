package dynamo

import (
	"fmt"
	"strings"
	"time"
)

// BuildPeriodUser builds a combined string of period and username required to identify an item of certain period and user.
func BuildPeriodUser(username, period string) *string {
	p := fmt.Sprintf("%s:%s", period, username)
	return &p
}

// BuildAmountKey builds the amount sort key, which is a combined string of the amount and the item ID.
// Example -> Input: 1234.56, abc123 -> Output: 000001234.56:abc123
func BuildAmountKey(amount float64, id string) string {
	amountStr := fmt.Sprintf("%0.2f", amount)
	return fmt.Sprintf("%012s:%s", amountStr, id)
}

// BuildNameKey builds the name sort key, which is a combined string of the name and the item ID.
func BuildNameKey(name, id string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(name), id)
}

// BuildCreatedDateEntityIDKey builds a variation of a created_date sort key, which is a combined string of the
// created_date and the item ID. This is required on scenarios where just using the created_date as the sort key is not
// enough to identify a single item.
func BuildCreatedDateEntityIDKey(createdDate time.Time, id string) string {
	return fmt.Sprintf("%s:%s", createdDate.Format(time.RFC3339Nano), id)
}
