package tasks_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type GetTasksResponseDTO []TaskDTO

// GetTasks godoc
// @Summary     Get tasks list
// @Description Get tasks list. Optional user_id filters by author. Total count is returned in X-Total-Count.
// @Tags        tasks
// @Produce     json
// @Param       user_id query int false "Author user ID"
// @Param       limit query int false "Limit"
// @Param       offset query int false "Offset"
// @Success     200 {array} TaskDTO "Tasks list"
// @Header      200 {integer} X-Total-Count "Total number of tasks"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [get]
func (h *TaskHttpHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Info("getting tasks")

	limit, offset, userId, err := getUserIdAndLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", err, core_errors.InvalidParameters),
			"failed to get user_id or limit or offset query params",
		)
		return
	}

	tasks, total, err := h.taskService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	tasksResponse := GetTasksResponseDTO(taskDtosFromDomains(tasks))

	rw.Header().Set("X-Total-Count", strconv.Itoa(total))
	responseHandler.JSONResponse(tasksResponse, http.StatusOK)
}

func getUserIdAndLimitOffsetQueryParams(req *http.Request) (limit *int, offset *int, userId *int, err error) {
	userId, err = core_http_request.GetIntQueryParams(req, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	if userId != nil && *userId <= 0 {
		return nil, nil, nil, fmt.Errorf("user_id must be greater than 0")
	}

	limit, err = core_http_request.GetIntQueryParams(req, "limit")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	if limit != nil && *limit <= 0 {
		return nil, nil, nil, fmt.Errorf("limit must be greater than 0")
	}

	offset, err = core_http_request.GetIntQueryParams(req, "offset")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	if offset != nil && *offset < 0 {
		return nil, nil, nil, fmt.Errorf("offset must be greater than or equal to 0")
	}

	return limit, offset, userId, nil
}
