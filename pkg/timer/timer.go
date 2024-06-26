package timer

import (
	"context"
	"time"

	"github.com/vuquang23/poseidon/pkg/logger"
)

func Start(ctx context.Context, task interface{}) func() {
	logger.Infof(ctx, "Start %v ...", task)

	start := time.Now()

	return func() {
		logger.Infof(ctx, "Finish %v in: %v", task, time.Since(start))
	}
}
