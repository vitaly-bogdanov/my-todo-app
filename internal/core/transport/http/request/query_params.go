package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
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

func GetDateQueryParams(req *http.Request, key string) (*time.Time, error) {
	queryParam := req.URL.Query().Get(key)
	if queryParam == "" {
		return nil, nil
	}

	layout := "2006-01-02"
	value, err := time.Parse(layout, queryParam)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid date %v: %w",
			queryParam,
			key,
			err,
			core_errors.BadRequest,
		)
	}

	return &value, nil
}
