package web_transport_http

import (
	"context"
	"net/http"

	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage(ctx context.Context) ([]byte, error)
	GetFile(ctx context.Context, name string) ([]byte, error)
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{webService: webService}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: h.GetMainPage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/assets/{file...}",
			Handler: h.GetAsset,
		},
	}
}
