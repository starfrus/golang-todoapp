package statistic_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/starfrus/golang-todoapp/internal/core/domain"
	core_http_server "github.com/starfrus/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	)(domain.Statistics, error)
}

func NewStatisticHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes()[]core_http_server.Route{
	return []core_http_server.Route{
		{
			Method: http.MethodGet,
			Path: "/statistics",
			Handler: h.GetStatistics,
		},
	}
}