package tasks_transport_http

import (
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type CreateTaskRequestDTO struct {
	Title        string  `json:"title"          validate:"required,min=1,max=100"`
	Description  *string `json:"description"    validate:"omitempty,min=1,max=1000"`
	AuthorUserID int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponseDTO TaskDTO

func (h *TaskHttpHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Info("creating task")

	var createTaskRequest CreateTaskRequestDTO
	if err := core_http_request.DecodeAndValidate(r, &createTaskRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	taskDomain := domainFromDTO(createTaskRequest)

	taskDomain, err := h.taskService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	createTaskResponse := CreateTaskResponseDTO(taskDtoFromDomain(taskDomain))

	responseHandler.JSONResponse(createTaskResponse, http.StatusCreated)
}

func domainFromDTO(dto CreateTaskRequestDTO) core_domain.Task {
	return core_domain.NewTaskUninitialized(dto.Title, dto.Description, dto.AuthorUserID)
}
