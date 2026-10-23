package users_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_repository_postgres_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type GetUserRequestDTO UserDTO

// GetUser godoc
// @Summary     Get user
// @Description Get user by ID
// @Tags        users
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} GetUserRequestDTO "User found"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [get]
func (h *UserHttpHandler) GetUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get user")

	userID, err := core_http_request.GetIntPathValue(req, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", err, core_errors.InvalidParameters),
			"failed to get user",
		)
		return
	}

	user, err := h.userService.GetUser(ctx, userID)
	if errors.Is(err, core_repository_postgres_pool.ErrNoRows) {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", err, core_errors.NotFound),
			"user not found",
		)
		return
	}
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	userDto := GetUserRequestDTO(userDtoFromDomain(user))

	log.Info("success")
	responseHandler.JSONResponse(userDto, http.StatusOK)
}
