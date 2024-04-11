package worker

import (
	"github.com/urfave/cli/v2"

	"github.com/vuquang23/poseidon/internal/pkg/config"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	tasksvc "github.com/vuquang23/poseidon/internal/pkg/service/task"
	"github.com/vuquang23/poseidon/internal/pkg/taskq/worker"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/eth"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
)

func RunWorker(c *cli.Context) error {
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

	// eth client
	ethClient, err := eth.NewClient(conf.Eth)
	if err != nil {
		return err
	}

	// repository
	poolRepo := poolrepo.New(db, asynqClient)
	txRepo := txrepo.New(db)

	// service
	taskSvc := tasksvc.New(conf.Service.Task, poolRepo, txRepo, ethClient, asynqClient)

	// worker
	w, err := worker.New(conf.Redis)
	if err != nil {
		return err
	}

	worker.RegisterHandlers(w, taskSvc)

	return w.Run()
}
