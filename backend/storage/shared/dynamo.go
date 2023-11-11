package shared

import "fmt"

// BuildPeriodUser builds a combined string of period and username required to identify an item of certain period and user.
func BuildPeriodUser(username, period string) *string {
	p := fmt.Sprintf("%s:%s", period, username)
	return &p
}
