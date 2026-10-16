package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	"golang.org/x/sync/errgroup"
)

func (s *TaskService) GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]core_domain.Task, int, error) {
	if limit != nil && *limit <= 0 {
		return nil, 0, fmt.Errorf("limit must be greater than 0")
	}

	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be greater than or equal to 0")
	}

	g, gctx := errgroup.WithContext(ctx)

	var tasks []core_domain.Task
	g.Go(func() error {
		var err error
		tasks, err = s.taskRepository.GetTasks(gctx, userId, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get tasks: %w", err)
		}
		return nil
	})

	var total int
	g.Go(func() error {
		var err error
		total, err = s.taskRepository.CountTasks(gctx, userId)
		if err != nil {
			return fmt.Errorf("failed to count tasks: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
