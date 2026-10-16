package tasks_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
)

type TaskHttpHandler struct {
	taskService TaskService
}

type TaskService interface {
	CreateTask(ctx context.Context, taskDomain core_domain.Task) (core_domain.Task, error)
	GetTasks(ctx context.Context, userId *int, limit *int, offset *int) ([]core_domain.Task, int, error)
	GetTask(ctx context.Context, taskId int) (core_domain.Task, error)
	DeleteTask(ctx context.Context, taskId int) error
	PatchTask(ctx context.Context, taskId int, patch core_domain.TaskPatch) (core_domain.Task, error)
}

func NewTaskHttpHandler(taskService TaskService) *TaskHttpHandler {
	return &TaskHttpHandler{taskService: taskService}
}

func (h *TaskHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
