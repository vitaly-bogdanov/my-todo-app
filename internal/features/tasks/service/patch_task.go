package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (s *TaskService) PatchTask(
	ctx context.Context,
	taskId int,
	patch core_domain.TaskPatch,
) (core_domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, taskId)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	if err := task.Apply(patch); err != nil {
		return core_domain.Task{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	patchedTask, err := s.taskRepository.PatchTask(ctx, taskId, task)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("failed to patch task: %w", err)
	}

	return patchedTask, nil
}
