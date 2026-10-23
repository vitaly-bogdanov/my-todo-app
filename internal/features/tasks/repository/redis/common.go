package tasks_repository_redis

import (
	"context"
	"fmt"
	"strconv"
)

func taskKey(taskID int) string {
	return fmt.Sprintf("task:%d", taskID)
}

func tasksListKey(userID *int) string {
	if userID == nil {
		return "tasks:all"
	}
	return fmt.Sprintf("tasks:%d", *userID)
}

func tasksListField(limit *int, offset *int) string {
	return fmt.Sprintf("%s:%s", pageValue(limit), pageValue(offset))
}

func pageValue(value *int) string {
	if value == nil {
		return "nil"
	}
	return strconv.Itoa(*value)
}

func (r *CachedRepository) invalidateTasks(ctx context.Context, userID int, taskID *int) {
	keys := []string{tasksListKey(nil), tasksListKey(&userID)}
	if taskID != nil {
		keys = append(keys, taskKey(*taskID))
	}

	if _, err := r.pool.Del(ctx, keys...).Result(); err != nil {
		logCacheError(ctx, "invalidate tasks cache", err)
	}
}
