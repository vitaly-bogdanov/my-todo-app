package core_errors

import "errors"

var (
	NotFound          = errors.New("not found")
	InvalidParameters = errors.New("invalid parameters")
	Conflict          = errors.New("conflict")
)
