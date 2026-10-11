package users_repository_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	if limit != nil && *limit > 0 {
		userModels = make([]UserModel, 0, *limit)
	}
	for rows.Next() {
		var user UserModel
		err := rows.Scan(
			&user.ID,
			&user.Version,
			&user.FullName,
			&user.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userModels = append(userModels, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return userDomainsFromModels(userModels), nil
}
