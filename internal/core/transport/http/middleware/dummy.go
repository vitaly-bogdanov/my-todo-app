package core_http_middleware

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
)

func Dummy() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			log.Info("dummy middleware before")
			next.ServeHTTP(w, r)
			log.Info("dummy middleware after")
		})
	}
}
