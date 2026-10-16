package tasks_service

import (
	"context"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type TaskService struct {
	taskRepository TaskRepository
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task core_domain.Task) (core_domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]core_domain.Task, error)
	CountTasks(ctx context.Context, userId *int) (int, error)
	GetTask(ctx context.Context, taskId int) (core_domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
	PatchTask(ctx context.Context, taskId int, taskDomain core_domain.Task) (core_domain.Task, error)
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}
