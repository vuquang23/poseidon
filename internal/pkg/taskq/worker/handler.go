package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	tasksvc "github.com/vuquang23/poseidon/internal/pkg/service/task"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/timer"
)

func RegisterHandlers(worker *Worker, taskSvc tasksvc.ITaskService) {
	worker.RegisterHandler(valueobject.TaskTypeHandlePoolCreated, HandlePoolCreated(taskSvc))
	worker.RegisterHandler(valueobject.TaskTypeScanTxs, ScanTxs(taskSvc))
	worker.RegisterHandler(valueobject.TaskTypeGetETHUSDTKline, GetETHUSDTKline(taskSvc))
	worker.RegisterHandler(valueobject.TaskTypeFinalizeTxs, FinalizeTxs(taskSvc))
}

func bindLoggerCtx(ctx context.Context, taskID, taskType string) context.Context {
	l := logger.WithFieldsNonContext(logger.Fields{"taskId": taskID, "taskType": taskType})
	return context.WithValue(ctx, logger.CtxLoggerKey, l)
}

func HandlePoolCreated(taskSvc tasksvc.ITaskService) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)
		ctx = bindLoggerCtx(ctx, taskID, valueobject.TaskTypeHandlePoolCreated)

		finish := timer.Start(ctx, taskID)
		defer finish()

		var payload valueobject.TaskHandlePoolCreatedPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return taskSvc.HandlePoolCreated(ctx, payload.PoolAddress)
	}
}

func ScanTxs(taskSvc tasksvc.ITaskService) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)
		ctx = bindLoggerCtx(ctx, taskID, valueobject.TaskTypeScanTxs)

		finish := timer.Start(ctx, taskID)
		defer finish()

		var payload valueobject.TaskScanTxsPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		err := taskSvc.ScanTxs(ctx, payload)
		if err != nil && err != tasksvc.ErrInvalidBlockRange {
			return err
		}

		return nil
	}
}

func GetETHUSDTKline(taskSvc tasksvc.ITaskService) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)
		ctx = bindLoggerCtx(ctx, taskID, valueobject.TaskTypeGetETHUSDTKline)

		finish := timer.Start(ctx, taskID)
		defer finish()

		var payload valueobject.TaskGetETHUSDTKlinePayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		return taskSvc.GetETHUSDTKline(ctx, payload)
	}
}

func FinalizeTxs(taskSvc tasksvc.ITaskService) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)
		ctx = bindLoggerCtx(ctx, taskID, valueobject.TaskTypeFinalizeTxs)

		finish := timer.Start(ctx, taskID)
		defer finish()

		var payload valueobject.TaskFinalizeTxsPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			logger.Error(ctx, err.Error())
			return err
		}

		err := taskSvc.FinalizeTxs(ctx, payload)
		if err != nil && err != tasksvc.ErrInvalidBlockRange {
			return err
		}

		return nil
	}
}
