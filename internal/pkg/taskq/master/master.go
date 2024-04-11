package master

import (
	"context"
	"time"

	"github.com/hibiken/asynq"

	"github.com/vuquang23/poseidon/internal/pkg/service/task"
	asynqpkg "github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/redis"
)

type Master struct {
	taskCfgProvider *periodicTaskCfgProvider
	taskManager     *asynq.PeriodicTaskManager
}

func New(cfg redis.Config, svc task.ITaskService) (*Master, error) {
	redisConnOpt, err := asynqpkg.GetAsynqRedisConnectionOption(cfg)
	if err != nil {
		return nil, err
	}

	provider := &periodicTaskCfgProvider{svc: svc}
	manager, err := asynq.NewPeriodicTaskManager(asynq.PeriodicTaskManagerOpts{
		RedisConnOpt:               redisConnOpt,
		PeriodicTaskConfigProvider: provider,
		SyncInterval:               15 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Master{
		taskCfgProvider: provider,
		taskManager:     manager,
	}, nil
}

func (m *Master) Run() error {
	return m.taskManager.Run()
}

type periodicTaskCfgProvider struct {
	svc task.ITaskService
}

func (p *periodicTaskCfgProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	return p.svc.GetPeriodicTaskConfigs(context.Background())
}
