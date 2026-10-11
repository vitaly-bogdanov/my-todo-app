package users_repository_postgres

import (
	"context"
	"fmt"
)

func (r *UserRepository) CountUsers(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM todoapp.users`).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return total, nil
}
