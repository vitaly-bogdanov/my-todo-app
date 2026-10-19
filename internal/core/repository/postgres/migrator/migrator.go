// Package core_repository_postgres_migrator применяет встроенные SQL-миграции при старте приложения.
package core_repository_postgres_migrator

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // драйвер pgx5://
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// migrationsTable — таблица версий golang-migrate. Схема указана явно: иначе драйвер берёт
// current_schema(), а она меняется, если имя пользователя БД совпадает со схемой из
// миграций (search_path = "$user", public). Тот же public.schema_migrations
// по умолчанию использует и CLI migrate из Makefile.
const migrationsTable = `"public"."schema_migrations"`

// Up применяет все ещё не применённые миграции из migrationsFS к базе databaseURL.
// databaseURL должен иметь схему pgx5://. Если новых миграций нет, Up ничего не делает
// и возвращает текущую версию схемы.
func Up(migrationsFS fs.FS, databaseURL string) (version uint, dirty bool, err error) {
	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return 0, false, fmt.Errorf("failed to open migrations source: %w", err)
	}

	u, err := url.Parse(databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to parse database url: %w", err)
	}
	q := u.Query()
	q.Set("x-migrations-table", migrationsTable)
	q.Set("x-migrations-table-quoted", "1")
	u.RawQuery = q.Encode()

	m, err := migrate.NewWithSourceInstance("iofs", source, u.String())
	if err != nil {
		return 0, false, fmt.Errorf("failed to init migrate: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if err == nil {
			err = errors.Join(srcErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, false, fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err = m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("failed to get migrations version: %w", err)
	}

	return version, dirty, nil
}
