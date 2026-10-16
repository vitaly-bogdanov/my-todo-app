package tasks_repository_postgres

import (
	"context"
	"errors"
	"fmt"
	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

func (r *TaskRepository) CreateTask(
	ctx context.Context,
	task core_domain.Task,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_repository_postgres_pool.ErrViolatesForeignKey) {
			return core_domain.Task{}, fmt.Errorf(
				"%v: user with id=%d not found: %w",
				err,
				task.AuthorUserID,
				core_errors.NotFound,
			)
		}
		return core_domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := core_domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)

	return taskDomain, nil
}
