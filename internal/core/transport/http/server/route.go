package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []core_http_middleware.Middleware
}

func (r *Route) WithMiddlewares() http.Handler {
	return core_http_middleware.ChainMiddlewares(
		r.Handler,
		r.Middlewares...,
	)
}
