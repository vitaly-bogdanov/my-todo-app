package users_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/utils"
)

func (h *UserHttpHandler) DeleteUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("delete user")

	userID, err := core_http_utils.GetIntPathValue(req, "id")
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
