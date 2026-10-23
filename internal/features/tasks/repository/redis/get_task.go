package tasks_repository_redis

import (
	"context"
	"encoding/json"
	"errors"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
)

func (r *CachedRepository) GetTask(ctx context.Context, taskID int) (core_domain.Task, error) {
	raw, err := r.pool.Get(ctx, taskKey(taskID)).Result()
	if err == nil {
		var model TaskModel
		if err := json.Unmarshal([]byte(raw), &model); err != nil {
			logCacheError(ctx, "decode task cache", err)
		} else {
			return taskDomainFromModel(model), nil
		}
	} else if !errors.Is(err, core_repository_redis_pool.NotFound) {
		logCacheError(ctx, "read task from cache", err)
	}

	task, err := r.mainRepository.GetTask(ctx, taskID)
	if err != nil {
		return core_domain.Task{}, err
	}

	body, err := json.Marshal(taskModelFromDomain(task))
	if err != nil {
		logCacheError(ctx, "encode task cache", err)
		return task, nil
	}
	if err := r.pool.Set(ctx, taskKey(task.ID), body, r.pool.TTL()).Err(); err != nil {
		logCacheError(ctx, "write task cache", err)
	}
	return task, nil
}
