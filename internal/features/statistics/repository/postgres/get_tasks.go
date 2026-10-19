package statistics_repository_postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) ([]core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT
			id,
			title,
			description,
			completed,
			created_at,
			completed_at
		FROM todoapp.tasks
	`)

	args := []any{}
	conditions := []string{}

	if userId != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, userId)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
