package core_repository_postgres_pool_pgx

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

type pgxRow struct {
	pgx.Row
}

func (r *pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core_repository_postgres_pool.ErrNoRows
	}

	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		if pgxErr.Code == "23503" {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_repository_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		core_repository_postgres_pool.ErrUnknown,
	)
}
