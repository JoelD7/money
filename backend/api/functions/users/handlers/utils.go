package handlers

import (
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"strconv"
)

func getRequestQueryParams(req *apigateway.Request) (string, int, error) {
	pageSizeParam := 0
	var err error

	if req.QueryStringParameters["page_size"] != "" {
		pageSizeParam, err = strconv.Atoi(req.QueryStringParameters["page_size"])
	}

	if err != nil || pageSizeParam < 0 {
		return "", 0, fmt.Errorf("%w: %v", models.ErrInvalidPageSize, err)
	}

	return req.QueryStringParameters["start_key"], pageSizeParam, nil
}
