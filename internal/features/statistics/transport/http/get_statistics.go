package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	core_logger "github.com/vitaly-bogdanov/my-todo-app/internal/core/logger"
	core_http_request "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/vitaly-bogdanov/my-todo-app/internal/core/transport/http/response"
)

type GetStatisticsResponseDTO struct {
	TasksCreated                 int      `json:"tasks_created"`
	TasksCompleted               int      `json:"tasks_completed"`
	TasksCompletedRate           *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletationTime *string  `json:"tasks_average_completation_time"`
}

// GetStatistics godoc
// @Summary     Get statistics
// @Description Get task statistics. from and to are dates in YYYY-MM-DD. If both are set, to must be after from.
// @Tags        statistics
// @Produce     json
// @Param       user_id query int false "Author user ID"
// @Param       from query string false "Start date, YYYY-MM-DD"
// @Param       to query string false "End date, YYYY-MM-DD"
// @Success     200 {object} GetStatisticsResponseDTO "Statistics"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /statistics [get]
func (h *StatisticsHttpHandler) GetStatistics(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	log.Info("get users list")

	userId, from, to, err := getUserIdFromToQueryParams(req)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id or from or to query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDtoFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDtoFromDomain(statistics core_domain.Statistics) GetStatisticsResponseDTO {
	var avgTime *string
	if statistics.TasksAverageCompletationTime != nil {
		duration := statistics.TasksAverageCompletationTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponseDTO{
		TasksCreated:                 statistics.TasksCreated,
		TasksCompleted:               statistics.TasksCompleted,
		TasksCompletedRate:           statistics.TasksCompletedRate,
		TasksAverageCompletationTime: avgTime,
	}
}

func getUserIdFromToQueryParams(req *http.Request) (userId *int, from *time.Time, to *time.Time, err error) {
	userId, err = core_http_request.GetIntQueryParams(req, "user_id")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err = core_http_request.GetDateQueryParams(req, "from")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err = core_http_request.GetDateQueryParams(req, "to")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userId, from, to, nil
}
