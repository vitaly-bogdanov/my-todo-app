package statistics_repository_postgres

import core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewStatisticsRepository(pool core_repository_postgres_pool.Pool) *StatisticsRepository {
	return &StatisticsRepository{pool: pool}
}
