package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
	"go.uber.org/zap"
)

// requestIDHeader — заголовок, по которому один HTTP-вызов связывается между сервисами.
const requestIDHeader = "X-Request-ID"

// RequestID возвращает middleware, которая проставляет requestIDHeader в запрос и в ответ.
// Если во входящем запросе заголовка нет, RequestID генерирует новый UUID.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger возвращает middleware, которая пишет запрос в log.
// Она читает requestIDHeader, чтобы строку лога можно было сопоставить с вызывающей стороной.
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			ctx := context.WithValue(r.Context(), core_logger.LoggerContextKey, log.With(
				zap.String("requestID", requestID),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Panic возвращает middleware, которая перехватывает панику в обработчике.
// Клиенту уходит ответ 500, текст паники пишется в лог из контекста.
// Логер туда заранее кладёт Logger.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicRecoveryResponse(r, "during request handling")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Trace возвращает middleware, которая пишет в лог начало и конец обработки запроса.
// В конце она добавляет длительность и статус ответа.
// Логер Trace берёт из контекста, его туда кладёт Logger.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rwWithStatusCode := core_http_response.NewResponseWriter(rw)

			before := time.Now().UTC()
			log.Debug(
				">>> incoming request",
				zap.Time("time", before),
			)

			next.ServeHTTP(rwWithStatusCode, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("statusCode", rwWithStatusCode.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}
