package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTO

func (h *TaskHttpHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Info("getting task")

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id path param")
		return
	}

	task, err := h.taskService.GetTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	taskResponse := GetTaskResponse(taskDtoFromDomain(task))

	responseHandler.JSONResponse(taskResponse, http.StatusOK)
}
