package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_repository_postgres_pool_pgx "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
	users_repository_postgres "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/service"
	users_transport_http "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to create logger:", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Info("creating connection pool...")

	pool, err := core_repository_postgres_pool_pgx.NewPool(
		ctx,
		core_repository_postgres_pool_pgx.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to create connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing features...", zap.String("features", "users"))
	usersRepository := users_repository_postgres.NewUserRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHttpHandler(usersService)

	logger.Debug("initializing HTTP server...")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIVersionRouter(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to start HTTP server: %w", zap.Error(err))
	}
}
