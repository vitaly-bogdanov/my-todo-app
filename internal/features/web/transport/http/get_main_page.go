package web_transport_http

import (
	"net/http"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get main page")

	html, err := h.webService.GetMainPage(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get main page")
		return
	}

	responseHandler.HtmlResponse(html)
}
