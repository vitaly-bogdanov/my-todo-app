package users_transport_http

import (
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

// CreateUserRequestDTO — тело запроса на создание пользователя.
type CreateUserRequestDTO struct {
	FullName    string  `json:"full_name"    validate:"required,min=1,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164,startswith=+"`
}

// CreateUserResponseDTO — тело ответа после создания пользователя.
type CreateUserResponseDTO UserDTO

// CreateUser обрабатывает POST /users и декодирует JSON-тело в CreateUserRequestDTO.
func (h *UserHttpHandler) CreateUser(rw http.ResponseWriter, req *http.Request) {

	// ... оптимизировать: ппробовать сделать валидотор-миддлвару
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Info("creating user")

	var createUserRequest CreateUserRequestDTO

	// попробовать вшить валидацию в само DTO как метод
	if err := core_http_request.DecodeAndValidate(req, &createUserRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userDomain := domainFromDTO(createUserRequest)

	userDomain, err := h.userService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	createUserResponse := CreateUserResponseDTO(userDtoFromDomain(userDomain))

	log.Info("success")
	responseHandler.JSONResponse(createUserResponse, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequestDTO) core_domain.User {
	return core_domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
