package statistic_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/starfrus/golang-todoapp/internal/core/domain"
	core_logger "github.com/starfrus/golang-todoapp/internal/core/logger"
	core_http_request "github.com/starfrus/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/starfrus/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int     `json:"tasks_created"`
	TasksCompleted             int     `json:"tasks_completed"`
	TasksCompletedRate         float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string `json:"tasks_average_completion_time"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandelr := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParam(r)
	if err != nil {
		responseHandelr.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	statistics, err := h.statisticService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandelr.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statistics)

	responseHandelr.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         *statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParam(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDParamKey    = "user_id"
		fromQueryParamKey = "from"
		toQueryParamKey   = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
