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
}

func bindLoggerCtx(ctx context.Context, taskID string) context.Context {
	l := logger.WithFieldsNonContext(logger.Fields{"taskId": taskID})
	return context.WithValue(ctx, logger.CtxLoggerKey, l)
}

func HandlePoolCreated(taskSvc tasksvc.ITaskService) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)
		ctx = bindLoggerCtx(ctx, taskID)

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
