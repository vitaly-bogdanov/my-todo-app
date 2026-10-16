package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (s *TaskService) GetTask(
	ctx context.Context,
	taskId int,
) (core_domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, taskId)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}
