package core_errors

import "errors"

var (
	InternalServerError = errors.New("internal server error")
	NotFound            = errors.New("not found")
	BadRequest          = errors.New("bad request")
	Unauthorized        = errors.New("unauthorized")
	Forbidden           = errors.New("forbidden")
	Conflict            = errors.New("conflict")
	TooManyRequests     = errors.New("too many requests")
	BadGateway          = errors.New("bad gateway")
	ServiceUnavailable  = errors.New("service unavailable")
)
