package tasks_repository_postgres

import (
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
)

type TaskRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewTaskRepository(
	pool core_repository_postgres_pool.Pool,
) *TaskRepository {
	return &TaskRepository{pool: pool}
}
