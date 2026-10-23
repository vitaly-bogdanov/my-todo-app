package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/vitaly-bogdanov/my-todo-app/docs"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      *Config
	log         *core_logger.Logger
	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middlewares ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      &config,
		log:         log,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		s.mux.Handle(pattern, route.WithMiddlewares())
	}
}

func (s *HTTPServer) RegisterAPIVersionRouter(routers ...*APIVersionRouter) {
	for _, router := range routers {
		handler := router.WithMiddlewares()
		for _, prefix := range []string{
			fmt.Sprintf("/api/%s", router.apiVersion),
			fmt.Sprintf("/%s", router.apiVersion),
		} {
			s.mux.Handle(
				prefix+"/",
				http.StripPrefix(prefix, handler),
			)
		}
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddlewares(s.mux, s.middlewares...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}
	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Info("starting HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
	case <-ctx.Done():
		s.log.Info("shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
		s.log.Info("HTTP server shut down")
	}

	return nil
}
