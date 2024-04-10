package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

func HandleTaskX() func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		return nil
	}
}
