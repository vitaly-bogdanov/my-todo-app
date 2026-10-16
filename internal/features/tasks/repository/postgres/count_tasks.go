package tasks_repository_postgres

import (
	"context"
	"fmt"
)

func (r *TaskRepository) CountTasks(ctx context.Context, userId *int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT COUNT(*) FROM todoapp.tasks`
	args := []any{}
	if userId != nil {
		query = `SELECT COUNT(*) FROM todoapp.tasks WHERE author_id = $1`
		args = append(args, userId)
	}

	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	return total, nil
}
