package web_transport_http

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetAsset(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get frontend asset")

	name := r.PathValue("file")
	body, err := h.webService.GetFile(ctx, path.Join("assets", name))
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get frontend asset")
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	responseHandler.FileResponse(body, contentType)
}
