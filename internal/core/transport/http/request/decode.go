package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

var requestValidator = validator.New()

type validateble interface {
	Validate() error
}

func DecodeAndValidate[T any](r *http.Request, dest T) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("failed to decode request body: %v: %w", err, core_errors.BadRequest)
	}

	var err error

	v, ok := any(dest).(validateble)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"failed to validate request body: %v: %w",
			err,
			core_errors.BadRequest,
		)
	}

	return nil
}
