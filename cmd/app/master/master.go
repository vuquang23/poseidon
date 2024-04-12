package master

import (
	"github.com/urfave/cli/v2"

	"github.com/vuquang23/poseidon/internal/pkg/config"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	tasksvc "github.com/vuquang23/poseidon/internal/pkg/service/task"
	"github.com/vuquang23/poseidon/internal/pkg/taskq/master"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
)

func RunMaster(c *cli.Context) error {
	conf := config.New()
	if err := conf.Load(c.String("config")); err != nil {
		return err
	}

	// logger
	_, err := logger.Init(conf.Log, logger.LoggerBackendZap)
	if err != nil {
		return err
	}

	// postgres
	db, err := postgres.New(conf.Postgres)
	if err != nil {
		return err
	}

	// asynq client
	asynqClient, err := asynq.NewClient(conf.Redis)
	if err != nil {
		return err
	}

	// repository
	poolRepo := poolrepo.New(db, asynqClient)

	// service
	taskSvc := tasksvc.New(conf.Service.Task, poolRepo, nil, nil, nil, asynqClient, nil)

	m, err := master.New(conf.Redis, taskSvc)
	if err != nil {
		return err
	}

	return m.Run()
}
