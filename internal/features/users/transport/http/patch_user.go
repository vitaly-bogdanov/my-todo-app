package users_transport_http

import (
	"fmt"
	"net/http"
	"regexp"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
	core_http_types "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/types"
)

var phoneNumberPattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

type PatchUserRequestDTO struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+1234567890"`
}

func (r *PatchUserRequestDTO) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("full name can't be null")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 1 || fullNameLen > 100 {
			return fmt.Errorf("full name must be between 1 and 100 characters")
		}
	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumber := *r.PhoneNumber.Value
		if !phoneNumberPattern.MatchString(phoneNumber) {
			return fmt.Errorf("invalid phone number")
		}
	}

	return nil
}

type PatchUserResponseDTO UserDTO

// 404 доделать

// PatchUser godoc
// @Summary     Patch user
// @Description Patch user by ID
// @Description
// @Description В теле можно передать:
// @Description - только phone_number
// @Description - только full_name
// @Description - phone_number и full_name
// @Description
// @Description full_name не может быть null. phone_number начинается с + и содержит от 8 до 15 цифр, первая цифра от 1 до 9.
// @Tags        users
// @Produce     json
// @Param       id path int true "User ID"
// @Param       request body PatchUserRequestDTO true "User patch"
// @Success     200 {object} PatchUserResponseDTO "User patched"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Not found"
// @Failure     409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [patch]
func (h *UserHttpHandler) PatchUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("patch user")

	userID, err := core_http_request.GetIntPathValue(req, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}

	var patchUserRequest PatchUserRequestDTO
	if err := core_http_request.DecodeAndValidate(req, &patchUserRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userPatch := userPatchFromRequest(patchUserRequest)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	userResponse := PatchUserResponseDTO(userDtoFromDomain(userDomain))

	log.Info("success")
	responseHandler.JSONResponse(userResponse, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequestDTO) core_domain.UserPatch {
	return core_domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
