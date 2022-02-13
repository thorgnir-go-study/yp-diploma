package accrualprocessor

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"time"
)

type Reader interface {
	GetTasksToRun(ctx context.Context) ([]*entity.ProcessingTask, error)
}

type Writer interface {
	CreateTask(ctx context.Context, task *entity.ProcessingTask) error
	CreateTasks(ctx context.Context, tasks []*entity.ProcessingTask) error
	SetTaskStatus(ctx context.Context, taskID entity.ID, status entity.ProcessingTaskStatus) error
	RescheduleTask(ctx context.Context, taskID entity.ID, nextRun time.Time) error
	CleanProcessedTasks(ctx context.Context) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Start() error
}
