package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIntQueryParams(req *http.Request, key string) (*int, error) {
	queryParam := req.URL.Query().Get(key)
	if queryParam == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(queryParam)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not an integer: %w",
			queryParam,
			key,
			err,
		)
	}

	return &value, nil
}
