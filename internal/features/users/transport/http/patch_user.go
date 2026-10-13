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

var phoneNumberPattern = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

type PatchUserRequestDTO struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" validate:"min=1,max=100"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" validate:"omitempty,e164,startswith=+"`
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
		phoneNumberLen := len([]rune(phoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("phone number must be between 10 and 15 characters")
		}

		if !phoneNumberPattern.MatchString(phoneNumber) {
			return fmt.Errorf("invalid phone number")
		}
	}

	return nil
}

type PatchUserResponseDTO UserDTO

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
