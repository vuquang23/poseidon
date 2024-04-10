package worker

import (
	"context"

	"github.com/hibiken/asynq"

	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	asynqpkg "github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/redis"
)

type Worker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func New(cfg redis.Config) (*Worker, error) {
	redisConnOpt, err := asynqpkg.GetAsynqRedisConnectionOption(cfg)
	if err != nil {
		return nil, err
	}

	server := asynq.NewServer(redisConnOpt, asynq.Config{})
	mux := asynq.NewServeMux()
	worker := &Worker{
		server: server,
		mux:    mux,
	}

	return worker, nil
}

func (w *Worker) Run() error {
	return w.server.Run(w.mux)
}

func (w *Worker) RegisterHandler(taskType valueobject.TaskType, handler func(ctx context.Context, t *asynq.Task) error) {
	w.mux.HandleFunc(string(taskType), handler)
}
