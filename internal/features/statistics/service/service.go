package statistics_service

import (
	"context"
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) ([]core_domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{statisticsRepository: statisticsRepository}
}
