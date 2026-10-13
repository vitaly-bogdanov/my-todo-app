package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIntPathValue(req *http.Request, key string) (int, error) {
	pathValue := req.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("path value is required")
	}

	value, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' by key='%s' is not an integer %v: %w",
			pathValue,
			key,
			err,
		)
	}

	return value, nil
}
