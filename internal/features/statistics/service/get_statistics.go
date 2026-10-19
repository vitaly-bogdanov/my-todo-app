package statistics_service

import (
	"context"
	"fmt"
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (core_domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return core_domain.Statistics{}, fmt.Errorf("to date is before from date: %w", core_errors.BadRequest)
		}
	}

	task, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return core_domain.Statistics{}, fmt.Errorf("failed to get tasks: %w", err)
	}

	statistics := calculateStatistics(task)

	return statistics, nil
}

func calculateStatistics(tasks []core_domain.Task) core_domain.Statistics {
	if len(tasks) == 0 {
		return core_domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	// TODO: fix this
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return core_domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)
}
