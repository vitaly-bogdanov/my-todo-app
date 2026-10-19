package core_domain

import "time"

type Statistics struct {
	TasksCreated                 int
	TasksCompleted               int
	TasksCompletedRate           *float64
	TasksAverageCompletationTime *time.Duration
}

func NewStatistics(
	TasksCreated int,
	TasksCompleted int,
	TasksCompletedRate *float64,
	TasksAverageCompletationTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:                 TasksCreated,
		TasksCompleted:               TasksCompleted,
		TasksCompletedRate:           TasksCompletedRate,
		TasksAverageCompletationTime: TasksAverageCompletationTime,
	}
}
