package tasks_repository_redis

import (
	"context"
	"encoding/json"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (r *CachedRepository) PatchTask(
	ctx context.Context,
	taskID int,
	task core_domain.Task,
) (core_domain.Task, error) {
	patched, err := r.mainRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return core_domain.Task{}, err
	}

	body, err := json.Marshal(taskModelFromDomain(patched))
	if err != nil {
		logCacheError(ctx, "encode task cache", err)
	} else if err := r.pool.Set(ctx, taskKey(patched.ID), body, r.pool.TTL()).Err(); err != nil {
		logCacheError(ctx, "write task cache", err)
	}
	r.invalidateTasks(ctx, patched.AuthorUserID, nil)
	return patched, nil
}
