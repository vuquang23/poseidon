package api

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"

	"github.com/vuquang23/poseidon/internal/pkg/api"
	"github.com/vuquang23/poseidon/internal/pkg/config"
	poolrepo "github.com/vuquang23/poseidon/internal/pkg/repository/pool"
	pricerepo "github.com/vuquang23/poseidon/internal/pkg/repository/price"
	txrepo "github.com/vuquang23/poseidon/internal/pkg/repository/tx"
	"github.com/vuquang23/poseidon/internal/pkg/server"
	poolsvc "github.com/vuquang23/poseidon/internal/pkg/service/pool"
	txsvc "github.com/vuquang23/poseidon/internal/pkg/service/tx"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/logger"
	"github.com/vuquang23/poseidon/pkg/postgres"
)

func RunAPI(c *cli.Context) error {
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

	// auto migration
	if err := postgres.MigrateUp(db, "file://./migration/postgres", 0); err != nil && err != migrate.ErrNoChange {
		return err
	}

	// repository
	poolRepo := poolrepo.New(db, asynqClient)
	txRepo := txrepo.New(db)
	priceRepo := pricerepo.New(db)

	// service
	poolSvc := poolsvc.New(poolRepo)
	txSvc := txsvc.New(txRepo, priceRepo)

	// server
	server := server.GinEngine(conf.Http, conf.Log, logger.LoggerBackendZap)

	api.RegisterRoutes(&conf, server, poolSvc, txSvc)

	return server.Run(conf.Http.BindAddress)
}
