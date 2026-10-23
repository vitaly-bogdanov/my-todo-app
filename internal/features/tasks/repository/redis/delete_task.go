package tasks_repository_redis

import (
	"context"
	"encoding/json"
	"errors"

	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
)

func (r *CachedRepository) DeleteTask(ctx context.Context, taskID int) error {
	authorUserID, ok := r.authorUserIDFromCache(ctx, taskID)
	if !ok {
		task, err := r.mainRepository.GetTask(ctx, taskID)
		if err != nil {
			return err
		}
		authorUserID = task.AuthorUserID
	}

	if err := r.mainRepository.DeleteTask(ctx, taskID); err != nil {
		return err
	}

	r.invalidateTasks(ctx, authorUserID, &taskID)
	return nil
}

func (r *CachedRepository) authorUserIDFromCache(ctx context.Context, taskID int) (int, bool) {
	raw, err := r.pool.Get(ctx, taskKey(taskID)).Result()
	if err != nil {
		if !errors.Is(err, core_repository_redis_pool.NotFound) {
			logCacheError(ctx, "read task from cache", err)
		}
		return 0, false
	}

	var model TaskModel
	if err := json.Unmarshal([]byte(raw), &model); err != nil {
		logCacheError(ctx, "decode task cache", err)
		return 0, false
	}

	return taskDomainFromModel(model).AuthorUserID, true
}
