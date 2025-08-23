package models

import (
	"fmt"
	"net/url"
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
	Active       bool
}

func (qp *QueryParameters) ParseAsURLValues(query *url.Values) {
	if qp.Period != "" {
		query.Add("period", qp.Period)
	}

	if qp.StartKey != "" {
		query.Add("start_key", qp.StartKey)
	}

	if qp.PageSize != 0 {
		query.Add("page_size", fmt.Sprint(qp.PageSize))
	}

	if qp.SavingGoalID != "" {
		query.Add("saving_goal_id", qp.SavingGoalID)
	}

	if qp.SortBy != "" {
		query.Add("sort_by", qp.SortBy)
	}

	if qp.SortType != "" {
		query.Add("sort_order", qp.SortType)
	}

	for _, category := range qp.Categories {
		query.Add("category", category)
	}
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
