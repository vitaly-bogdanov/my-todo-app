package tasks_service

import (
	"context"
	"fmt"
)

func (s *TaskService) DeleteTask(ctx context.Context, taskId int) error {
	err := s.taskRepository.DeleteTask(ctx, taskId)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
