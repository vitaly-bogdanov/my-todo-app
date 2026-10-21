package users_transport_http

import (
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type CreateUserRequestDTO struct {
	FullName    string  `json:"full_name"    validate:"required,min=1,max=100"      example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164,startswith=+" example:"+1234567890"`
}

type CreateUserResponseDTO UserDTO

// CreateUser   godoc
// @Summary     Create a new user
// @Description Create a new user with the given full name and phone number
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body CreateUserRequestDTO true "User creation request"
// @Success     201 {object} CreateUserResponseDTO "User created successfully"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users [post]
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

	responseHandler.JSONResponse(createUserResponse, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequestDTO) core_domain.User {
	return core_domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
