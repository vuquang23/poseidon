package task

import (
	"context"

	"github.com/hibiken/asynq"

	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
)

type ITaskService interface {
	GetPeriodicTaskConfigs(ctx context.Context) ([]*asynq.PeriodicTaskConfig, error)

	HandlePoolCreated(ctx context.Context, poolAddress string) error
	ScanTxs(ctx context.Context, payload valueobject.TaskScanTxsPayload) error
	GetETHUSDTKline(ctx context.Context, payload valueobject.TaskGetETHUSDTKlinePayload) error
}
