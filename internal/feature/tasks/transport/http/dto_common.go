package tasks_transport_http

import (
	"time"

	"github.com/starfrus/golang-todoapp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"competed_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func taskDTOFormDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomain(tasks []domain.Task) []TaskDTOResponse {
	dts := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dts[i] = taskDTOFormDomain(task)
	}

	return dts
}
