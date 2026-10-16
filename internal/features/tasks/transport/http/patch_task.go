package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/types"
)

type PatchTaskRequestDTO struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequestDTO) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can't be null")
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("Title must be between 1 and 100 characters")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("Description must be between 1 and 1000 characters")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed can't be null")
		}
	}

	return nil
}

type PatchTaskResponseDTO TaskDTO

func (h *TaskHttpHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id path param")
		return
	}

	var patchTaskRequest PatchTaskRequestDTO
	if err := core_http_request.DecodeAndValidate(r, &patchTaskRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	taskPatch := taskPatchFromRequest(patchTaskRequest)

	task, err := h.taskService.PatchTask(ctx, taskId, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponseDTO(taskDtoFromDomain(task))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequestDTO) core_domain.TaskPatch {
	return core_domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
