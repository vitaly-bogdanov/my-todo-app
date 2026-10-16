package tasks_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

func (r *TaskRepository) PatchTask(
	ctx context.Context,
	taskId int,
	task core_domain.Task,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.tasks
		SET 
			title=$1,
			description=$2,
			completed=$3,
			completed_at=$4,
			version=version+1

		WHERE id=$5 AND version=$6

		RETURNING
			id,
			version,
			title,
			description,
			completed,
			completed_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		taskId,
		task.Version,
	)

	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CompletedAt,
	); err != nil {
		if errors.Is(err, core_repository_postgres_pool.ErrNoRows) {
			return core_domain.Task{}, fmt.Errorf("task with id=%d cuncurrently accesed: %w", taskId, core_errors.Conflict)
		}

		return core_domain.Task{}, fmt.Errorf("failed to scan task: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
