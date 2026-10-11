package users_repository_postgres

import (
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

type UserRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewUserRepository(pool core_repository_postgres_pool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}
