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

// GetStatusCodeOrPanic возвращает статус, записанный через WriteHeader.
// Если WriteHeader ещё не вызывали, GetStatusCodeOrPanic паникует.
func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("status code is not initialized")
	}
	return rw.statusCode
}
