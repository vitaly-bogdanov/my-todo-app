package users_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

// DeleteUser godoc
// @Summary     Delete user
// @Description Delete user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     204 "No content"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [delete]
func (h *UserHttpHandler) DeleteUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("delete user")

	userID, err := core_http_request.GetIntPathValue(req, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID")
		return
	}

	err = h.userService.DeleteUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	log.Info("success")
	responseHandler.NoContentResponse()
}
