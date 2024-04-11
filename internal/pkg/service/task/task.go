package task

import (
	"context"

	"github.com/hibiken/asynq"
)

type ITaskService interface {
	GetPeriodicTaskConfigs(ctx context.Context) ([]*asynq.PeriodicTaskConfig, error)

	HandlePoolCreated(ctx context.Context, poolAddress string) error
}
