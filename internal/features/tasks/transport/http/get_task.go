package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTO

// GetTask godoc
// @Summary     Get task
// @Description Get task by ID
// @Tags        tasks
// @Produce     json
// @Param       id path int true "Task ID"
// @Success     200 {object} GetTaskResponse "Task found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks/{id} [get]
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
