package tasks_service

import (
	"context"
	"fmt"
	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (s *TaskService) CreateTask(
	ctx context.Context,
	taskDomain core_domain.Task,
) (core_domain.Task, error) {
	if err := taskDomain.Validate(); err != nil {
		return core_domain.Task{}, fmt.Errorf("invalid task domain: %w", err)
	}

	taskDomain, err := s.taskRepository.CreateTask(ctx, taskDomain)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return taskDomain, nil
}