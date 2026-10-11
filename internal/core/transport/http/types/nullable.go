package core_http_types

import (
	"encoding/json"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type Nullable[T any] struct {
	core_domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Value = &value
	return nil
}

func (n *Nullable[T]) ToDomain() core_domain.Nullable[T] {
	return core_domain.Nullable[T]{
		Set:   n.Set,
		Value: n.Value,
	}
}
