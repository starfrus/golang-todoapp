package domain

import "time"

type Statistics struct{
	TasksCreated int
	TasksCompleted int
	TasksCompletedRate *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistic(
	tasksCreated int,
	tasksCompelted int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
)Statistics{
	return Statistics{
		TasksCreated: tasksCreated,
		TasksCompleted: tasksCompelted,
		TasksCompletedRate: tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}

