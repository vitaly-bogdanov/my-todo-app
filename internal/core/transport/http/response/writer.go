package core_http_response

import "net/http"

// StatusCodeUninitialized означает, что обработчик ещё не записал статус ответа.
const StatusCodeUninitialized = -1

// ResponseWriter запоминает HTTP-статус, который обработчик записал в ответ.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter оборачивает w и ставит статус в StatusCodeUninitialized.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, statusCode: StatusCodeUninitialized}
}

// WriteHeader записывает statusCode в ответ и запоминает его.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// GetStatusCode возвращает статус, записанный через WriteHeader.
// Если WriteHeader ещё не вызывали, GetStatusCode паникует.
func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}
	return rw.statusCode
}
