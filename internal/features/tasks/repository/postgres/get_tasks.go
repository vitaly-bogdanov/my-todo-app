package tasks_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (r *TaskRepository) GetTasks(
	ctx context.Context,
	userId *int,
	limit *int,
	offset *int,
) ([]core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_id
		FROM todoapp.tasks
		%s
		ORDER BY created_at ASC
		LIMIT $1
		OFFSET $2
	`

	args := []any{limit, offset}

	if userId != nil {
		query = fmt.Sprintf(query, "WHERE author_id = $3")
		args = append(args, userId)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskModel
	for rows.Next() {
		var task TaskModel
		err := rows.Scan(
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
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	taskDomains := taskDomainsFromModels(tasks)

	return taskDomains, nil
}
