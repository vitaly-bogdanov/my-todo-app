package users_repository_postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (r *UserRepository) PatchUser(
	ctx context.Context,
	userID int,
	user core_domain.User,
) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.users
		SET 
			full_name=$1,
			phone_number=$2,
			version=version+1
		WHERE id=$3 AND version=$4
		RETURNING 
			id,
			version,
			full_name,
			phone_number
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		userID,
		user.Version,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf(
				"user with id=%d cuncurrently accesed: %w",
				userID,
				core_errors.NotFound,
			)
		}

		return core_domain.User{}, fmt.Errorf("failed to scan user: %w", err)
	}

	userDomain := core_domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
