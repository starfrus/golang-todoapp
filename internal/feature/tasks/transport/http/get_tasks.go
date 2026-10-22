package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/starfrus/golang-todoapp/internal/core/logger"
	core_http_request "github.com/starfrus/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/starfrus/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	user_id, limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `userID/limit/offset` query param")
		return
	}

	tasksDomain, err := h.tasksService.GetTasks(ctx, user_id, limit, offset)
	if err != nil{
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomain(tasksDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDParamKey      = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `userID` query paran: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
