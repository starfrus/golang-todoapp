package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/starfrus/golang-todoapp/internal/core/domain"
	core_errors "github.com/starfrus/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvaligArgument,
			)
		}
	}

	// 1.tasks := get tasks
	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}
	// 2. stat := calcStat()

	statistics := calcStatistics(tasks)

	// return stat
	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistic(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)
	var totalCompletedDuration time.Duration
	tasksCompleted := 0
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++

		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistic(
		tasksCreated, 
		tasksCompleted, 
		&tasksCompletedRate, 
		tasksAverageCompletionTime,
	)
}
