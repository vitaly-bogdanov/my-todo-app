package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

func (h *TaskHttpHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Info("delete task")

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id path param")
		return
	}

	err = h.taskService.DeleteTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)

}
