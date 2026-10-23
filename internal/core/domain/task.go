package core_domain

import (
	"fmt"
	"time"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}

	if t.CompletedAt == nil {
		return nil
	}

	duration := t.CompletedAt.Sub(t.CreatedAt)

	return &duration
}

// Validate проверяет заголовок, описание и согласованность статуса завершения задачи.
func (t *Task) Validate() error {
	// Title обязателен: от 1 до 100 символов.
	titleLength := len([]rune(t.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("title must be between 1 and 100 characters: %d: %w", titleLength, core_errors.InvalidParameters)
	}

	// Description необязателен; если задан, длина от 1 до 1000 символов.
	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))
		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf("description must be between 1 and 1000 characters: %d: %w", descriptionLength, core_errors.InvalidParameters)
		}
	}

	if t.Completed {
		// Завершённая задача должна иметь CompletedAt.
		if t.CompletedAt == nil {
			return fmt.Errorf("completed at is required when task is completed: %w", core_errors.InvalidParameters)
		}
		// CompletedAt не может быть раньше CreatedAt.
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("completed at must be after created at: %w", core_errors.InvalidParameters)
		}
	} else {
		// Незавершённая задача не должна иметь CompletedAt.
		if t.CompletedAt != nil {
			return fmt.Errorf("completed at must be nil when task is not completed: %w", core_errors.InvalidParameters)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (t *TaskPatch) Validate() error {
	if t.Title.Set && t.Title.Value == nil {
		return fmt.Errorf("Title can't be patched to null: %w", core_errors.InvalidParameters)
	}

	if t.Completed.Set && t.Completed.Value == nil {
		return fmt.Errorf("Completed can't be patched to null: %w", core_errors.InvalidParameters)
	}

	return nil
}

func (t *Task) Apply(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("failed to validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
		if tmp.Completed {
			now := time.Now()
			tmp.CompletedAt = &now
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate task: %w", err)
	}

	*t = tmp

	return nil
}
