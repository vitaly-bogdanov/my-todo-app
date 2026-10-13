package core_repository_postgres_pool_pgx

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

// ------------------------------ pgx.Rows ------------------------------

type pgxRows struct {
	pgx.Rows
}

// ------------------------------ pgx.CommandTag -----------------------

type pgxCommandTag struct {
	pgconn.CommandTag
}

// ------------------------------ pgx.Row ------------------------------

type pgxRow struct {
	pgx.Row
}

func (r *pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_repository_postgres_pool.ErrNoRows
		}
		return err
	}
	return nil
}
