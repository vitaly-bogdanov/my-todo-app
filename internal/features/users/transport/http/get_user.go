package users_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/utils"
)

type GetUserRequest UserDTO

func (h *UserHttpHandler) GetUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get user")

	userID, err := core_http_utils.GetIntPathValue(req, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", err, core_errors.BadRequest),
			"failed to get user",
		)
		return
	}

	user, err := h.userService.GetUser(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
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

	userDto := userDtoFromDomain(user)

	log.Info("success")
	responseHandler.JSONResponse(userDto, http.StatusOK)
}
