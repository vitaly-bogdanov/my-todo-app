package users_transport_http

import (
	"context"
	"net/http"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
)

type UserHttpHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]core_domain.User, int, error)
	GetUser(ctx context.Context, userID int) (core_domain.User, error)
	DeleteUser(ctx context.Context, userID int) error
	PatchUser(ctx context.Context, userID int, userPatch core_domain.UserPatch) (core_domain.User, error)
}

func NewUserHttpHandler(userService UserService) *UserHttpHandler {
	return &UserHttpHandler{userService: userService}
}

func (h *UserHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			Middlewares: []core_http_middleware.Middleware{
				core_http_middleware.Dummy(),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
