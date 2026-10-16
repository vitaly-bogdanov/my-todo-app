package tasks_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

func (r *TaskRepository) GetTask(
	ctx context.Context,
	taskId int,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	  SELECT id, version, title, description, completed, created_at, completed_at, author_id
		FROM todoapp.tasks
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, taskId)

	var task TaskModel
	err := row.Scan(
		&task.ID,
		&task.Version,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_repository_postgres_pool.ErrNoRows) {
			return core_domain.Task{}, fmt.Errorf("task with id=%d not found %w: %w", taskId, err, core_errors.NotFound)
		}
		return core_domain.Task{}, fmt.Errorf("failed to get task id=%d: %w", taskId, err)
	}

	taskDomain := taskDomainFromModel(task)

	return taskDomain, nil
}
