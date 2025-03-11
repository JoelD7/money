package models

import (
	"strconv"
	"strings"
)

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"

	SortParamCreatedDate SortParam = "created_date"
	SortParamAmount      SortParam = "amount"
	SortParamName        SortParam = "name"
	SortParamDeadline    SortParam = "deadline"
	SortParamTarget      SortParam = "target"
)

type SortOrder string
type SortParam string

type QueryParameters struct {
	Categories   []string
	Period       string
	StartKey     string
	PageSize     int
	SortBy       string
	SortType     string
	SavingGoalID string
}

func (qp *QueryParameters) ToURLParams() string {
	urlParams := make([]string, 0)

	for _, category := range qp.Categories {
		urlParams = append(urlParams, "category="+category)
	}

	if qp.Period != "" {
		urlParams = append(urlParams, "period="+qp.Period)
	}

	if qp.PageSize != 0 {
		urlParams = append(urlParams, "page_size="+strconv.Itoa(qp.PageSize))
	}

	if qp.SortBy != "" {
		urlParams = append(urlParams, "sort_by="+qp.SortBy)
	}

	if qp.SortType != "" {
		urlParams = append(urlParams, "sort_type="+qp.SortType)
	}

	if qp.SavingGoalID != "" {
		urlParams = append(urlParams, "saving_goal_id="+qp.SavingGoalID)
	}

	if qp.StartKey != "" {
		urlParams = append(urlParams, "start_key="+qp.StartKey)
	}

	return strings.Join(urlParams, "&")
}
