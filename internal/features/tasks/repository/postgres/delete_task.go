package tasks_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (r *TaskRepository) DeleteTask(ctx context.Context, taskId int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.tasks
		WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id=%d not found: %w", taskId, core_errors.NotFound)
	}

	return nil
}
