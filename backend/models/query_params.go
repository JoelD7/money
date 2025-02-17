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
type SortParam string

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"

	SortParamCreatedDate SortParam = "created_date"
	SortParamAmount      SortParam = "amount"
	SortParamName        SortParam = "name"
	SortParamDeadline    SortParam = "deadline"
	SortParamTarget      SortParam = "target"
)
