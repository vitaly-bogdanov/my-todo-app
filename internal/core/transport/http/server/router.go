package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersionV1 = ApiVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion  ApiVersion
	middlewares []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion ApiVersion,
	middlewares ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  apiVersion,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddlewares())
	}
}

func (r *APIVersionRouter) WithMiddlewares() http.Handler {
	return core_http_middleware.ChainMiddlewares(
		r,
		r.middlewares...,
	)
}
