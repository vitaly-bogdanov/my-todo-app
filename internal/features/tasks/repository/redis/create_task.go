package tasks_repository_redis

import (
	"context"
	"encoding/json"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (r *CachedRepository) CreateTask(ctx context.Context, task core_domain.Task) (core_domain.Task, error) {
	created, err := r.mainRepository.CreateTask(ctx, task)
	if err != nil {
		return core_domain.Task{}, err
	}

	body, err := json.Marshal(taskModelFromDomain(created))
	if err != nil {
		logCacheError(ctx, "encode task cache", err)
	} else if err := r.pool.Set(ctx, taskKey(created.ID), body, r.pool.TTL()).Err(); err != nil {
		logCacheError(ctx, "write task cache", err)
	}
	r.invalidateTasks(ctx, created.AuthorUserID, nil)
	return created, nil
}
