package users_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (r *UserRepository) GetUser(ctx context.Context, userID int) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		WHERE id = $1;
	`

	var userModel UserModel
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_domain.User{}, fmt.Errorf("user not found: %w", core_errors.NotFound)
	}
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return core_domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	), nil
}
