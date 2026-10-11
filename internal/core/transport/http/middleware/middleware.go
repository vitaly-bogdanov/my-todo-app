package core_http_middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddlewares(httpHandler http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return httpHandler
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		httpHandler = middlewares[i](httpHandler)
	}

	return httpHandler
}
