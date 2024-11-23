package apigateway

type QueryParameters struct {
	Categories   []string
	Period       string
	StartKey     string
	PageSize     int
	SortBy       string
	SortType     string
	SavingGoalID string
}
