package users_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (r *UserRepository) DeleteUser(
		ctx context.Context,
		userID int,
	) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.users WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("failed to get rows affected: %w", core_errors.NotFound)
	}

	return nil
}
