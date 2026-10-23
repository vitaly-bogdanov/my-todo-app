package tasks_repository_redis

import (
	"context"
	"encoding/json"
	"errors"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
)

func (r *CachedRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]core_domain.Task, error) {
	raw, err := r.pool.HGet(ctx, tasksListKey(userID), tasksListField(limit, offset)).Result()
	if err == nil {
		var models []TaskModel
		if err := json.Unmarshal([]byte(raw), &models); err != nil {
			logCacheError(ctx, "decode tasks cache", err)
		} else {
			tasks := make([]core_domain.Task, len(models))
			for i, model := range models {
				tasks[i] = taskDomainFromModel(model)
			}
			return tasks, nil
		}
	} else if !errors.Is(err, core_repository_redis_pool.NotFound) {
		logCacheError(ctx, "read tasks from cache", err)
	}

	tasks, err := r.mainRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	models := make([]TaskModel, len(tasks))
	for i, task := range tasks {
		models[i] = taskModelFromDomain(task)
	}
	body, err := json.Marshal(models)
	if err != nil {
		logCacheError(ctx, "encode tasks cache", err)
		return tasks, nil
	}
	key := tasksListKey(userID)
	if _, err := r.pool.HSet(ctx, key, tasksListField(limit, offset), body).Result(); err != nil {
		logCacheError(ctx, "write tasks cache", err)
		return tasks, nil
	}
	if err := r.pool.Expire(ctx, key, r.pool.TTL()).Err(); err != nil {
		logCacheError(ctx, "expire tasks cache", err)
	}
	return tasks, nil
}

func (r *CachedRepository) CountTasks(ctx context.Context, userID *int) (int, error) {
	return r.mainRepository.CountTasks(ctx, userID)
}
