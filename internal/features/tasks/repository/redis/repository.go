package tasks_repository_redis

import (
	"context"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
	"go.uber.org/zap"
)

type taskRepository interface {
	CreateTask(ctx context.Context, task core_domain.Task) (core_domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]core_domain.Task, error)
	CountTasks(ctx context.Context, userId *int) (int, error)
	GetTask(ctx context.Context, taskId int) (core_domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
	PatchTask(ctx context.Context, taskId int, taskDomain core_domain.Task) (core_domain.Task, error)
}

type CachedRepository struct {
	pool           core_repository_redis_pool.Pool
	mainRepository taskRepository
}

func NewCachedRepository(
	pool core_repository_redis_pool.Pool,
	mainRepository taskRepository,
) *CachedRepository {
	return &CachedRepository{pool: pool, mainRepository: mainRepository}
}

func logCacheError(ctx context.Context, message string, err error) {
	logger, ok := ctx.Value(core_logger.LoggerContextKey).(*core_logger.Logger)
	if !ok {
		return
	}
	logger.Error(message, zap.Error(err))
}
