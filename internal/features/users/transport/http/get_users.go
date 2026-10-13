package users_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

// GetUsersResponseDTO — тело ответа после получения списка пользователей.
type GetUsersResponseDTO []UserDTO

// GetUsers обрабатывает GET /users и возвращает список пользователей.
func (h *UserHttpHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get users list")

	limit, offset, err := getLimitOffsetQueryParams(req)
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("%w: %w", err, core_errors.BadRequest),
			"get limit offset query params",
		)
		return
	}

	users, total, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	usersResponse := GetUsersResponseDTO(userDtosFromDomains(users))

	rw.Header().Set("X-Total-Count", strconv.Itoa(total))
	log.Info("success")
	responseHandler.JSONResponse(usersResponse, http.StatusOK)
}

func getLimitOffsetQueryParams(req *http.Request) (*int, *int, error) {
	limit, err := core_http_request.GetIntQueryParams(req, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	if limit != nil && *limit <= 0 {
		return nil, nil, fmt.Errorf("limit must be greater than 0")
	}

	offset, err := core_http_request.GetIntQueryParams(req, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	if offset != nil && *offset < 0 {
		return nil, nil, fmt.Errorf("offset must be greater than or equal to 0")
	}

	return limit, offset, nil
}
