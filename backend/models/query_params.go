package models

type QueryParameters struct {
	Categories   []string
	Period       string
	StartKey     string
	PageSize     int
	SortBy       string
	SortType     string
	SavingGoalID string
}

type SortOrder string

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"
)
