package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_http_server "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/server"
)

type StatisticsHttpHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) (core_domain.Statistics, error)
}

func NewStatisticsHttpHandler(
	statisticsService StatisticsService,
) *StatisticsHttpHandler {
	return &StatisticsHttpHandler{statisticsService: statisticsService}
}

func (h *StatisticsHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
