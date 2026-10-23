package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/vitaly-bogdanov/my-todo-app/docs"
	core_config "github.com/vitaly-bogdanov/my-todo-app/internal/core/config"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_repository_postgres_migrator "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/migrator"
	core_repository_postgres_pool_pgx "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/postgres/pool/pgx"
	core_repository_redis_pool_goredis "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool/goredis"
	core_http_middleware "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
	statistics_repository_postgres "github.com/vitaly-bogdanov/my-todo-app/internal/features/statistics/repository/postgres"
	statistics_service "github.com/vitaly-bogdanov/my-todo-app/internal/features/statistics/service"
	statistics_transport_http "github.com/vitaly-bogdanov/my-todo-app/internal/features/statistics/transport/http"
	tasks_repository_postgres "github.com/vitaly-bogdanov/my-todo-app/internal/features/tasks/repository/postgres"
	tasks_repository_redis "github.com/vitaly-bogdanov/my-todo-app/internal/features/tasks/repository/redis"
	tasks_service "github.com/vitaly-bogdanov/my-todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/vitaly-bogdanov/my-todo-app/internal/features/tasks/transport/http"
	users_repository_postgres "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/service"
	users_transport_http "github.com/vitaly-bogdanov/my-todo-app/internal/features/users/transport/http"
	web_repository_file_system "github.com/vitaly-bogdanov/my-todo-app/internal/features/web/repository/file_system"
	web_service "github.com/vitaly-bogdanov/my-todo-app/internal/features/web/service"
	web_transport_http "github.com/vitaly-bogdanov/my-todo-app/internal/features/web/transport/http"
	"github.com/vitaly-bogdanov/my-todo-app/migrations"
	"go.uber.org/zap"
)

// @title       Todo App API
// @version     1.0
// @description Todo App API
// @host        localhost:8080
// @BasePath    /api/v1
// @schemes     http
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to create logger:", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Info("setting time zone...", zap.String("time zone", time.Local.String()))
	postgresConfig := core_repository_postgres_pool_pgx.NewConfigMust()

	logger.Info("applying database migrations...")
	migrationsVersion, migrationsDirty, err := core_repository_postgres_migrator.Up(
		migrations.FS,
		postgresConfig.URL("pgx5"),
	)
	if err != nil {
		logger.Fatal("failed to apply database migrations", zap.Error(err))
	}
	logger.Info(
		"database migrations applied",
		zap.Uint("version", migrationsVersion),
		zap.Bool("dirty", migrationsDirty),
	)

	logger.Info("creating connection pool...")

	pool, err := core_repository_postgres_pool_pgx.NewPool(
		ctx,
		postgresConfig,
	)
	if err != nil {
		logger.Fatal("failed to create connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Info("creating redis pool...")
	redisPool, err := core_repository_redis_pool_goredis.NewPool(
		ctx,
		core_repository_redis_pool_goredis.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to create redis pool", zap.Error(err))
	}
	defer redisPool.Close()

	logger.Debug("initializing features...", zap.String("features", "users"))
	usersRepository := users_repository_postgres.NewUserRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHttpHandler(usersService)

	logger.Debug("initializing features...", zap.String("features", "tasks"))
	tasksRepository := tasks_repository_postgres.NewTaskRepository(pool)
	cachedTasksRepository := tasks_repository_redis.NewCachedRepository(redisPool, tasksRepository)
	tasksService := tasks_service.NewTaskService(cachedTasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTaskHttpHandler(tasksService)

	logger.Debug("initializing features...", zap.String("features", "statistics"))
	statisticsRepository := statistics_repository_postgres.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHttpHandler(statisticsService)

	logger.Debug("initializing features...", zap.String("features", "web"))
	webRepository := web_repository_file_system.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server...")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(
			httpConfig.CORSAllowedOrigins,
			httpConfig.CORSAllowedMethods,
			httpConfig.CORSAllowedHeaders,
		),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)

	httpServer.RegisterAPIVersionRouter(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to start HTTP server: %w", zap.Error(err))
	}
}
